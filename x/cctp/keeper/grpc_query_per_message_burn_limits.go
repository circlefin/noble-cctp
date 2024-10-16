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

package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PerMessageBurnLimit(c context.Context, req *types.QueryGetPerMessageBurnLimitRequest) (*types.QueryGetPerMessageBurnLimitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetPerMessageBurnLimit(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPerMessageBurnLimitResponse{BurnLimit: val}, nil
}

func (k Keeper) PerMessageBurnLimits(c context.Context, req *types.QueryAllPerMessageBurnLimitsRequest) (*types.QueryAllPerMessageBurnLimitsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var perMessageBurnLimits []types.PerMessageBurnLimit
	ctx := sdk.UnwrapSDKContext(c)

	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	perMessageBurnLimitsStore := prefix.NewStore(adapter, types.KeyPrefix(types.PerMessageBurnLimitKeyPrefix))

	pageRes, err := query.Paginate(perMessageBurnLimitsStore, req.Pagination, func(key []byte, value []byte) error {
		var perMessageBurnLimit types.PerMessageBurnLimit
		if err := k.cdc.Unmarshal(value, &perMessageBurnLimit); err != nil {
			return err
		}

		perMessageBurnLimits = append(perMessageBurnLimits, perMessageBurnLimit)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPerMessageBurnLimitsResponse{BurnLimits: perMessageBurnLimits, Pagination: pageRes}, nil
}
