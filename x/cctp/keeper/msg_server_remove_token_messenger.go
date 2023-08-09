package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) RemoveTokenMessenger(goCtx context.Context, msg *types.MsgRemoveTokenMessenger) (*types.MsgRemoveTokenMessengerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, found := k.GetAuthority(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAuthorityNotSet, "authority is not set")
	}

	if owner.Address != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot remove token messengers")
	}

	existingTokenMessenger, found := k.GetTokenMessenger(ctx, msg.DomainId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrTokenMessengerNotFound, "no token messenger was found for this domain")
	}

	k.DeleteTokenMessenger(ctx, msg.DomainId)

	event := types.RemoteTokenMessengerRemoved{
		Domain:         msg.DomainId,
		TokenMessenger: []byte(strings.ToLower(existingTokenMessenger.Address)),
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgRemoveTokenMessengerResponse{}, err
}
