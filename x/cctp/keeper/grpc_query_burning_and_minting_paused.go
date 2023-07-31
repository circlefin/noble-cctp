package keeper

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) BurningAndMintingPaused(c context.Context, req *types.QueryGetBurningAndMintingPausedRequest) (*types.QueryGetBurningAndMintingPausedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetBurningAndMintingPaused(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetBurningAndMintingPausedResponse{Paused: val}, nil
}
