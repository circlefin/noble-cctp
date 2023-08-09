package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) UnlinkTokenPair(goCtx context.Context, msg *types.MsgUnlinkTokenPair) (*types.MsgUnlinkTokenPairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, found := k.GetAuthority(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrAuthorityNotSet, "authority is not set")
	}

	if owner.Address != msg.From {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "this message sender cannot unlink token pairs")
	}

	tokenPair, found := k.GetTokenPair(ctx, msg.RemoteDomain, strings.ToLower(msg.RemoteToken))
	if !found {
		return nil, sdkerrors.Wrap(types.ErrTokenPairNotFound, "token pair doesn't exist in store")
	}

	k.DeleteTokenPair(ctx, msg.RemoteDomain, strings.ToLower(msg.RemoteToken))

	event := types.TokenPairUnlinked{
		LocalToken:   tokenPair.LocalToken,
		RemoteDomain: tokenPair.RemoteDomain,
		RemoteToken:  tokenPair.RemoteToken,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)
	return &types.MsgUnlinkTokenPairResponse{}, err
}
