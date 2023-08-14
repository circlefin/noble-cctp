package keeper

import (
	"context"
	"strings"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) AddRemoteTokenMessenger(goCtx context.Context, msg *types.MsgAddRemoteTokenMessenger) (*types.MsgAddRemoteTokenMessengerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot add remote token messengers")
	}

	_, found := k.GetRemoteTokenMessenger(ctx, msg.DomainId)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrRemoteTokenMessengerAlreadyFound, "a remote token messenger for this domain already exists")
	}

	if msg.Address == "" {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	}

	newRemoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: msg.DomainId,
		Address:  msg.Address,
	}
	k.SetRemoteTokenMessenger(ctx, newRemoteTokenMessenger)

	event := types.RemoteTokenMessengerAdded{
		Domain:               msg.DomainId,
		RemoteTokenMessenger: []byte(strings.ToLower(msg.Address)),
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgAddRemoteTokenMessengerResponse{}, err
}
