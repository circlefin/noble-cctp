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

func (k Keeper) TokenMessenger(c context.Context, req *types.QueryGetTokenMessengerRequest) (*types.QueryGetTokenMessengerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTokenMessenger(ctx, req.DomainId)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTokenMessengerResponse{TokenMessenger: val}, nil
}

func (k Keeper) TokenMessengers(c context.Context, req *types.QueryAllTokenMessengersRequest) (*types.QueryAllTokenMessengersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tokenMessengers []types.TokenMessenger
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	tokenMessengersStore := prefix.NewStore(store, types.KeyPrefix(types.TokenMessengerKeyPrefix))

	pageRes, err := query.Paginate(tokenMessengersStore, req.Pagination, func(key []byte, value []byte) error {
		var tokenMessenger types.TokenMessenger
		if err := k.cdc.Unmarshal(value, &tokenMessenger); err != nil {
			return err
		}

		tokenMessengers = append(tokenMessengers, tokenMessenger)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTokenMessengersResponse{TokenMessengers: tokenMessengers, Pagination: pageRes}, nil
}
