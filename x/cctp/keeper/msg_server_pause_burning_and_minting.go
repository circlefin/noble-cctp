package keeper

import (
	"context"
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) PauseBurningAndMinting(goCtx context.Context, msg *types.MsgPauseBurningAndMinting) (*types.MsgPauseBurningAndMintingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, found := k.GetAuthority(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAuthorityNotSet, "authority is not set")
	}

	if owner.Address != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot pause burning and minting")
	}

	paused := types.BurningAndMintingPaused{
		Paused: true,
	}
	k.SetBurningAndMintingPaused(ctx, paused)

	event := types.BurningAndMintingPausedEvent{}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgPauseBurningAndMintingResponse{}, err

}
