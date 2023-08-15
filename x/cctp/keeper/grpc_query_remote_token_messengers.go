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

func (k Keeper) RemoteTokenMessenger(c context.Context, req *types.QueryRemoteTokenMessengerRequest) (*types.QueryRemoteTokenMessengerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetRemoteTokenMessenger(ctx, req.DomainId)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryRemoteTokenMessengerResponse{RemoteTokenMessenger: val}, nil
}

func (k Keeper) RemoteTokenMessengers(c context.Context, req *types.QueryRemoteTokenMessengersRequest) (*types.QueryRemoteTokenMessengersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var remoteTokenMessengers []types.RemoteTokenMessenger
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	remoteTokenMessengersStore := prefix.NewStore(store, types.KeyPrefix(types.RemoteTokenMessengerKeyPrefix))

	pageRes, err := query.Paginate(remoteTokenMessengersStore, req.Pagination, func(key []byte, value []byte) error {
		var remoteTokenMessenger types.RemoteTokenMessenger
		if err := k.cdc.Unmarshal(value, &remoteTokenMessenger); err != nil {
			return err
		}

		remoteTokenMessengers = append(remoteTokenMessengers, remoteTokenMessenger)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRemoteTokenMessengersResponse{RemoteTokenMessengers: remoteTokenMessengers, Pagination: pageRes}, nil
}
