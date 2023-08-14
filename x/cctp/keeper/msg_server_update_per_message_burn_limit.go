package keeper

import (
	"context"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) UpdatePerMessageBurnLimit(goCtx context.Context, msg *types.MsgUpdatePerMessageBurnLimit) (*types.MsgUpdatePerMessageBurnLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenController := k.GetTokenController(ctx)
	if tokenController != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot update the per message burn limit")
	}

	newPerMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  strings.ToLower(msg.Denom),
		Amount: msg.Amount,
	}
	k.SetPerMessageBurnLimit(ctx, newPerMessageBurnLimit)

	err := ctx.EventManager().EmitTypedEvent(msg)

	return &types.MsgUpdatePerMessageBurnLimitResponse{}, err
}
