package keeper

import (
	"context"
	"strings"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AddTokenMessenger(goCtx context.Context, msg *types.MsgAddTokenMessenger) (*types.MsgAddTokenMessengerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, found := k.GetAuthority(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAuthorityNotSet, "authority is not set")
	}

	if owner.Address != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot add token messengers")
	}

	_, found = k.GetTokenMessenger(ctx, msg.DomainId)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrTokenMessengerAlreadyFound, "a token messenger for this domain already exists")
	}

	if msg.Address == "" {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	}

	newTokenMessenger := types.TokenMessenger{
		DomainId: msg.DomainId,
		Address:  msg.Address,
	}
	k.SetTokenMessenger(ctx, newTokenMessenger)

	event := types.RemoteTokenMessengerAdded{
		Domain:         msg.DomainId,
		TokenMessenger: []byte(strings.ToLower(msg.Address)),
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgAddTokenMessengerResponse{}, err
}
