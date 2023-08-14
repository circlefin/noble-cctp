package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) DisableAttester(goCtx context.Context, msg *types.MsgDisableAttester) (*types.MsgDisableAttesterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	attesterManager := k.GetAttesterManager(ctx)
	if attesterManager != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot disable attesters")
	}

	_, found := k.GetAttester(ctx, string(msg.Attester))
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDisableAttester, "attester not found")
	}

	// disallow removing public key if there is only 1 active public key
	storedAttesters := k.GetAllAttesters(ctx)
	if len(storedAttesters) == 1 {
		return nil, sdkerrors.Wrap(types.ErrDisableAttester, "cannot disable the last attester")
	}

	// disallow removing public key if it causes the n in m/n multisig to fall below m (threshold # of signers)
	signatureThreshold, found := k.GetSignatureThreshold(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDisableAttester, "signature threshold not set")
	}

	if uint32(len(storedAttesters)) <= signatureThreshold.Amount {
		return nil, sdkerrors.Wrap(types.ErrDisableAttester, "signature threshold is too low")
	}

	k.DeleteAttester(ctx, string(msg.Attester))

	event := types.AttesterDisabled{
		Attester: string(msg.Attester),
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgDisableAttesterResponse{}, err
}
