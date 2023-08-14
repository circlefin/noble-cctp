package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) UnpauseBurningAndMinting(goCtx context.Context, msg *types.MsgUnpauseBurningAndMinting) (*types.MsgUnpauseBurningAndMintingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pauser := k.GetPauser(ctx)
	if pauser != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot unpause burning and minting")
	}

	paused := types.BurningAndMintingPaused{
		Paused: false,
	}
	k.SetBurningAndMintingPaused(ctx, paused)

	event := types.BurningAndMintingUnpausedEvent{}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUnpauseBurningAndMintingResponse{}, err
}
