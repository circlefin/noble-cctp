package keeper

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AcceptOwner(goCtx context.Context, msg *types.MsgAcceptOwner) (*types.MsgAcceptOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	currentOwner := k.GetOwner(ctx)
	pendingOwner, found := k.GetPendingOwner(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "pending owner is not set")
	}

	if pendingOwner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "you are not the pending owner")
	}

	k.SetOwner(ctx, pendingOwner)
	k.DeletePendingOwner(ctx)

	event := types.OwnerUpdated{
		PreviousOwner: currentOwner,
		NewOwner:      pendingOwner,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgAcceptOwnerResponse{}, err
}
