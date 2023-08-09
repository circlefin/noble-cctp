package keeper

import (
	"bytes"
	"context"
	sdkerrors "cosmossdk.io/errors"
	"encoding/hex"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (k msgServer) ReplaceDepositForBurn(goCtx context.Context, msg *types.MsgReplaceDepositForBurn) (*types.MsgReplaceDepositForBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	paused, found := k.GetBurningAndMintingPaused(ctx)
	if found && paused.Paused {
		return nil, sdkerrors.Wrap(types.ErrDepositForBurn, "burning and minting are paused")
	}

	// verify and parse original originalMessage
	if len(msg.OriginalMessage) < types.MessageBodyIndex {
		return nil, sdkerrors.Wrap(types.ErrDepositForBurn, "invalid message: too short")
	}
	originalMessage := DecodeMessage(msg.OriginalMessage)

	// verify and parse BurnMessage
	if len(originalMessage.MessageBody) != types.BurnMessageLen {
		return nil, sdkerrors.Wrap(types.ErrDepositForBurn, "burn message body is not the correct length")
	}
	burnMessage := DecodeBurnMessage(originalMessage.MessageBody)

	// validate originalMessage sender is the same as this message sender
	if msg.From != string(originalMessage.Sender) {
		return nil, sdkerrors.Wrap(types.ErrDepositForBurn, "invalid sender for message")
	}

	// validate new mint recipient
	emptyByteArr := make([]byte, types.MintRecipientLen)
	if bytes.Equal(emptyByteArr, msg.NewMintRecipient) {
		return nil, sdkerrors.Wrap(types.ErrDepositForBurn, "mint recipient must be nonzero")
	}

	newMessageBody := types.BurnMessage{
		Version:       burnMessage.Version,
		BurnToken:     burnMessage.BurnToken,
		MintRecipient: msg.NewMintRecipient,
		Amount:        burnMessage.Amount,
		MessageSender: burnMessage.MessageSender,
	}

	newMessageBodyBytes, err := EncodeBurnMessage(newMessageBody)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrParsingBurnMessage, "error parsing burn message")
	}

	replaceMessage := types.MsgReplaceMessage{
		From:                 msg.From,
		OriginalMessage:      msg.OriginalMessage,
		OriginalAttestation:  msg.OriginalAttestation,
		NewMessageBody:       newMessageBodyBytes,
		NewDestinationCaller: msg.NewDestinationCaller,
	}
	_, err = k.ReplaceMessage(goCtx, &replaceMessage)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "error calling replace message")
	}

	event := types.DepositForBurn{
		Nonce:                     originalMessage.Nonce,
		BurnToken:                 hex.EncodeToString(crypto.Keccak256(burnMessage.BurnToken)),
		Amount:                    burnMessage.Amount,
		Depositor:                 msg.From,
		MintRecipient:             msg.NewMintRecipient,
		DestinationDomain:         originalMessage.DestinationDomain,
		DestinationTokenMessenger: originalMessage.Recipient,
		DestinationCaller:         msg.NewDestinationCaller,
	}

	err = ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgReplaceDepositForBurnResponse{}, err
}
