package keeper

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateOwner(goCtx context.Context, msg *types.MsgUpdateOwner) (*types.MsgUpdateOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	currentOwner := k.GetOwner(ctx)
	if currentOwner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot update the authority")
	}

	k.SetOwner(ctx, msg.NewOwner)

	event := types.OwnerUpdated{
		PreviousOwner: currentOwner,
		NewOwner:      msg.NewOwner,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUpdateOwnerResponse{}, err
}
