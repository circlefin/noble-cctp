/*
 * Copyright (c) 2023, © Circle Internet Financial, LTD.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func TestAttesterQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNAttesters(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetAttesterRequest
		response *types.QueryGetAttesterResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAttesterRequest{
				Attester: "Attester0",
			},
			response: &types.QueryGetAttesterResponse{Attester: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAttesterRequest{
				Attester: "Attester1",
			},
			response: &types.QueryGetAttesterResponse{Attester: msgs[1]},
		},
		{
			desc: "AttesterNotFound",
			request: &types.QueryGetAttesterRequest{
				Attester: "nothing",
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Attester(wctx, tc.request)
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

func TestAttesterQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNAttesters(keeper, ctx, 5)
	Attester := make([]types.Attester, len(msgs))
	copy(Attester, msgs)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAttestersRequest {
		return &types.QueryAllAttestersRequest{
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
		for i := 0; i < len(Attester); i += step {
			resp, err := keeper.Attesters(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Attesters), step)
			require.Subset(t,
				nullify.Fill(Attester),
				nullify.Fill(resp.Attesters),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(Attester); i += step {
			resp, err := keeper.Attesters(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Attesters), step)
			require.Subset(t,
				nullify.Fill(Attester),
				nullify.Fill(resp.Attesters),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.Attesters(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(Attester), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(Attester),
			nullify.Fill(resp.Attesters),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.Attesters(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
	t.Run("PaginateError", func(t *testing.T) {
		_, err := keeper.Attesters(wctx, request([]byte("key"), 1, 0, true))
		require.Contains(t, err.Error(), "invalid request, either offset or key is expected, got both")
	})
}

func TestAttesterQueryPaginatedInvalidState(t *testing.T) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keeper, ctx := keepertest.CctpKeeperWithKey(t, storeKey)

	store := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.AttesterKeyPrefix))
	store.Set(types.KeyPrefix(string(types.AttesterKey([]byte("attester")))), []byte("invalid"))

	goCtx := sdk.WrapSDKContext(ctx)
	_, err := keeper.Attesters(goCtx, &types.QueryAllAttestersRequest{})

	parsedErr, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, parsedErr.Code(), codes.Internal)
}
