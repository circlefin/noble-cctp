package keeper

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateSignatureThreshold(goCtx context.Context, msg *types.MsgUpdateSignatureThreshold) (*types.MsgUpdateSignatureThresholdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	attesterManager := k.GetAttesterManager(ctx)
	if attesterManager != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot update the authority")
	}

	if msg.Amount == 0 {
		return nil, sdkerrors.Wrapf(types.ErrUpdateSignatureThreshold, "invalid signature threshold")
	}

	existingSignatureThreshold, _ := k.GetSignatureThreshold(ctx)
	if msg.Amount == existingSignatureThreshold.Amount {
		return nil, sdkerrors.Wrapf(types.ErrUpdateSignatureThreshold, "signature threshold already set to this value")
	}

	// new signature threshold cannot be greater than the number of stored public keys
	attesters := k.GetAllAttesters(ctx)
	if msg.Amount > uint32(len(attesters)) {
		return nil, sdkerrors.Wrapf(types.ErrUpdateSignatureThreshold, "new signature threshold is too high")
	}

	newSignatureThreshold := types.SignatureThreshold{
		Amount: msg.Amount,
	}

	k.SetSignatureThreshold(ctx, newSignatureThreshold)

	event := types.SignatureThresholdUpdated{
		OldSignatureThreshold: uint64(existingSignatureThreshold.Amount),
		NewSignatureThreshold: uint64(newSignatureThreshold.Amount),
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUpdateSignatureThresholdResponse{}, err
}
