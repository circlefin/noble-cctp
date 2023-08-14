package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) UpdateMaxMessageBodySize(goCtx context.Context, msg *types.MsgUpdateMaxMessageBodySize) (*types.MsgUpdateMaxMessageBodySizeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot update the max message body size")
	}

	newMaxMessageBodySize := types.MaxMessageBodySize{
		Amount: msg.MessageSize,
	}
	k.SetMaxMessageBodySize(ctx, newMaxMessageBodySize)

	event := types.MaxMessageBodySizeUpdated{
		NewMaxMessageBodySize: msg.MessageSize,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUpdateMaxMessageBodySizeResponse{}, err
}
