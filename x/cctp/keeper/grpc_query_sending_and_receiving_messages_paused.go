package keeper

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SendingAndReceivingMessagesPaused(c context.Context, req *types.QueryGetSendingAndReceivingMessagesPausedRequest) (*types.QueryGetSendingAndReceivingMessagesPausedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetSendingAndReceivingMessagesPaused(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetSendingAndReceivingMessagesPausedResponse{Paused: val}, nil
}
