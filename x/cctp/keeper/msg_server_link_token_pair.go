package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) LinkTokenPair(goCtx context.Context, msg *types.MsgLinkTokenPair) (*types.MsgLinkTokenPairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, found := k.GetAuthority(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAuthorityNotSet, "authority is not set")
	}

	if owner.Address != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot link token pairs")
	}

	// check whether there already exists a mapping for this remote domain/token
	_, found = k.GetTokenPair(ctx, msg.RemoteDomain, strings.ToLower(msg.RemoteToken))
	if found {
		return nil, sdkerrors.Wrapf(
			types.ErrTokenPairAlreadyFound,
			"local token for this remote domain + remote token mapping already exists in store")
	}

	newTokenPair := types.TokenPair{
		RemoteDomain: msg.RemoteDomain,
		RemoteToken:  strings.ToLower(msg.RemoteToken),
		LocalToken:   strings.ToLower(msg.LocalToken),
	}

	k.SetTokenPair(ctx, newTokenPair)

	event := types.TokenPairLinked{
		LocalToken:   newTokenPair.LocalToken,
		RemoteDomain: newTokenPair.RemoteDomain,
		RemoteToken:  strings.ToLower(newTokenPair.RemoteToken),
	}
	err := ctx.EventManager().EmitTypedEvent(&event)
	return &types.MsgLinkTokenPairResponse{}, err
}
