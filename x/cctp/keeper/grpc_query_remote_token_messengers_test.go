// Copyright 2024 Circle Internet Group, Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package keeper_test

import (
	"testing"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRemoteTokenMessengerQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper()
	msgs := createNRemoteTokenMessengers(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryRemoteTokenMessengerRequest
		response *types.QueryRemoteTokenMessengerResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryRemoteTokenMessengerRequest{
				DomainId: 0,
			},
			response: &types.QueryRemoteTokenMessengerResponse{RemoteTokenMessenger: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryRemoteTokenMessengerRequest{
				DomainId: 1,
			},
			response: &types.QueryRemoteTokenMessengerResponse{RemoteTokenMessenger: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryRemoteTokenMessengerRequest{
				DomainId: 123,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.RemoteTokenMessenger(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestRemoteTokenMessengerQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper()
	msgs := createNRemoteTokenMessengers(keeper, ctx, 5)
	RemoteTokenMessenger := make([]types.RemoteTokenMessenger, len(msgs))
	copy(RemoteTokenMessenger, msgs)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryRemoteTokenMessengersRequest {
		return &types.QueryRemoteTokenMessengersRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(RemoteTokenMessenger); i += step {
			resp, err := keeper.RemoteTokenMessengers(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RemoteTokenMessengers), step)
			require.Subset(t,
				nullify.Fill(RemoteTokenMessenger),
				nullify.Fill(resp.RemoteTokenMessengers),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(RemoteTokenMessenger); i += step {
			resp, err := keeper.RemoteTokenMessengers(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RemoteTokenMessengers), step)
			require.Subset(t,
				nullify.Fill(RemoteTokenMessenger),
				nullify.Fill(resp.RemoteTokenMessengers),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.RemoteTokenMessengers(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(RemoteTokenMessenger), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(RemoteTokenMessenger),
			nullify.Fill(resp.RemoteTokenMessengers),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.RemoteTokenMessengers(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
	t.Run("PaginateError", func(t *testing.T) {
		_, err := keeper.RemoteTokenMessengers(ctx, request([]byte("key"), 1, 0, true))
		require.Contains(t, err.Error(), "invalid request, either offset or key is expected, got both")
	})
}

func TestRemoteTokenMessengerQueryPaginatedInvalidState(t *testing.T) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	keeper, ctx := keepertest.CctpKeeperWithKey(storeKey)

	adapter := runtime.KVStoreAdapter(runtime.NewKVStoreService(storeKey).OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.RemoteTokenMessengerKeyPrefix))
	store.Set(types.RemoteTokenMessengerKey(0), []byte("invalid"))

	_, err := keeper.RemoteTokenMessengers(ctx, &types.QueryRemoteTokenMessengersRequest{})

	parsedErr, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, parsedErr.Code(), codes.Internal)
}
