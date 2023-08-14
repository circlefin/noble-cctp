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

func (k Keeper) UsedNonce(c context.Context, req *types.QueryGetUsedNonceRequest) (*types.QueryGetUsedNonceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	nonce := types.Nonce{
		SourceDomain: req.SourceDomain,
		Nonce:        req.Nonce,
	}
	found := k.GetUsedNonce(ctx, nonce)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetUsedNonceResponse{Nonce: nonce}, nil
}

func (k Keeper) UsedNonces(c context.Context, req *types.QueryAllUsedNoncesRequest) (*types.QueryAllUsedNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var usedNonces []types.Nonce
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	usedNonceStore := prefix.NewStore(store, types.KeyPrefix(types.UsedNonceKeyPrefix))

	pageRes, err := query.Paginate(usedNonceStore, req.Pagination, func(key []byte, value []byte) error {
		var usedNonce types.Nonce
		if err := k.cdc.Unmarshal(value, &usedNonce); err != nil {
			return err
		}

		usedNonces = append(usedNonces, usedNonce)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllUsedNoncesResponse{UsedNonces: usedNonces, Pagination: pageRes}, nil
}
