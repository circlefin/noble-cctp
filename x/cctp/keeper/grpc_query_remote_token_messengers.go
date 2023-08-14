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
