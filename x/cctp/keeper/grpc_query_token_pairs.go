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
package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TokenPair(c context.Context, req *types.QueryGetTokenPairRequest) (*types.QueryGetTokenPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTokenPairHex(ctx, req.RemoteDomain, req.RemoteToken)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTokenPairResponse{Pair: val}, nil
}

func (k Keeper) TokenPairs(c context.Context, req *types.QueryAllTokenPairsRequest) (*types.QueryAllTokenPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tokenPairs []types.TokenPair
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	TokenPairsStore := prefix.NewStore(store, types.KeyPrefix(types.TokenPairKeyPrefix))

	pageRes, err := query.Paginate(TokenPairsStore, req.Pagination, func(key []byte, value []byte) error {
		var tokenPair types.TokenPair
		if err := k.cdc.Unmarshal(value, &tokenPair); err != nil {
			return err
		}

		tokenPairs = append(tokenPairs, tokenPair)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTokenPairsResponse{TokenPairs: tokenPairs, Pagination: pageRes}, nil
}
