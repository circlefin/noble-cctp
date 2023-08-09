package keeper

import (
	"bytes"
	"context"

	"strings"

	sdkerrors "cosmossdk.io/errors"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	fiattokenfactorytypes "github.com/strangelove-ventures/noble/x/fiattokenfactory/types"
)

func (k msgServer) ReceiveMessage(goCtx context.Context, msg *types.MsgReceiveMessage) (*types.MsgReceiveMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	paused, found := k.GetSendingAndReceivingMessagesPaused(ctx)
	if found && paused.Paused {
		return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "sending and receiving messages are paused")
	}

	// Validate each signature in the attestation
	publicKeys := k.GetAllAttesters(ctx)
	if len(publicKeys) == 0 {
		return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "no attesters found")
	}

	signatureThreshold, found := k.GetSignatureThreshold(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "signature threshold not found")
	}

	verified, err := VerifyAttestationSignatures(msg.Message, msg.Attestation, publicKeys, signatureThreshold.Amount)
	if err != nil || !verified {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "unable to verify signatures")
	}

	// validate message format
	if len(msg.Message) < types.MessageBodyIndex {
		return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "invalid message: too short")
	}

	message := DecodeMessage(msg.Message)

	// validate domain
	if message.DestinationDomain != types.NobleDomainId {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "incorrect destination domain: %d", message.DestinationDomain)
	}

	// validate destination caller
	emptyByteArr := make([]byte, types.DestinationCallerLen)
	if !bytes.Equal(message.DestinationCaller, emptyByteArr) && string(message.DestinationCaller) != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "incorrect destination caller: %s, sender: %s", message.DestinationCaller, msg.From)
	}

	// validate version
	if message.Version != types.NobleMessageVersion {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "incorrect message version. expected: %d, found: %d", types.NobleMessageVersion, message.Version)
	}

	// validate nonce is available
	// note: we use the domain/nonce combo instead of a hash
	usedNonce := types.Nonce{SourceDomain: message.SourceDomain, Nonce: message.Nonce}
	found = k.GetUsedNonce(ctx, usedNonce)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "nonce already used")
	}

	// mark nonce as used
	k.SetUsedNonce(ctx, usedNonce)

	// verify and parse BurnMessage
	burnMessageIsValid := len(message.MessageBody) == types.BurnMessageLen
	burnMessage := DecodeBurnMessage(message.MessageBody)

	if burnMessage.Version != types.MessageBodyVersion {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "invalid message body version")
	}

	if burnMessageIsValid { // then mint
		// look up Noble mint token from corresponding source domain/token
		tokenPair, found := k.GetTokenPair(ctx, message.SourceDomain, strings.ToLower(string(burnMessage.BurnToken)))
		if !found {
			return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "corresponding noble mint token not found")
		}

		msgMint := fiattokenfactorytypes.MsgMint{
			From:    msg.From,
			Address: string(burnMessage.MintRecipient),
			Amount: sdk.Coin{
				Denom:  strings.ToLower(tokenPair.LocalToken),
				Amount: sdk.NewIntFromBigInt(burnMessage.Amount.BigInt()),
			},
		}

		_, err = k.fiattokenfactory.Mint(goCtx, &msgMint)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "Error during minting")
		}

		mintEvent := types.MintAndWithdraw{
			MintRecipient: string(burnMessage.MintRecipient),
			Amount:        burnMessage.Amount,
			MintToken:     strings.ToLower(tokenPair.LocalToken),
		}
		err = ctx.EventManager().EmitTypedEvent(&mintEvent)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "Error emitting mint event")
		}
	}

	// on failure to decode, nil err from handleMessage
	if err := k.router.HandleMessage(ctx, msg.Message); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrHandleMessage, "Error in handleMessage")
	}

	event := types.MessageReceived{
		Caller:       msg.From,
		SourceDomain: message.SourceDomain,
		Nonce:        message.Nonce,
		Sender:       message.Sender,
		MessageBody:  message.MessageBody,
	}
	err = ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgReceiveMessageResponse{Success: true}, err
}
