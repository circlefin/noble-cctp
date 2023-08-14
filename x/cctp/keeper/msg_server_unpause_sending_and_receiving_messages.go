package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) UnpauseSendingAndReceivingMessages(goCtx context.Context, msg *types.MsgUnpauseSendingAndReceivingMessages) (*types.MsgUnpauseSendingAndReceivingMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pauser := k.GetPauser(ctx)
	if pauser != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot unpause sending and receiving messages")
	}

	paused := types.SendingAndReceivingMessagesPaused{
		Paused: false,
	}
	k.SetSendingAndReceivingMessagesPaused(ctx, paused)

	event := types.SendingAndReceivingUnpausedEvent{}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUnpauseSendingAndReceivingMessagesResponse{}, err
}
