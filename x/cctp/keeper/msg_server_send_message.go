package keeper

import (
	"bytes"
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) SendMessage(goCtx context.Context, msg *types.MsgSendMessage) (*types.MsgSendMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nonce := k.ReserveAndIncrementNonce(ctx)

	err := k.sendMessage(
		ctx,
		msg.DestinationDomain,
		msg.Recipient,
		make([]byte, types.DestinationCallerLen),
		[]byte(msg.From),
		nonce.Nonce,
		msg.MessageBody)

	return &types.MsgSendMessageResponse{Nonce: nonce.Nonce}, err
}

func (k msgServer) sendMessage(
	ctx sdk.Context,
	destinationDomain uint32,
	recipient []byte,
	destinationCaller []byte,
	messageSender []byte,
	nonce uint64,
	messageBody []byte,
) error {
	paused, found := k.GetSendingAndReceivingMessagesPaused(ctx)
	if found && paused.Paused {
		return sdkerrors.Wrap(types.ErrSendMessage, "sending and receiving messages is paused")
	}

	// check if message body is too long, ignore if max length not found
	max, found := k.GetMaxMessageBodySize(ctx)
	if found && uint64(len(messageBody)) > max.Amount {
		return sdkerrors.Wrap(types.ErrSendMessage, "message body exceeds max size")
	}

	emptyByteArr := make([]byte, len(recipient))
	if len(recipient) == 0 || bytes.Equal(recipient, emptyByteArr) {
		return sdkerrors.Wrap(types.ErrSendMessage, "recipient must not be nonzero")
	}

	// serialize message
	message := types.Message{
		Version:           types.MessageBodyVersion,
		SourceDomain:      types.NobleDomainId,
		DestinationDomain: destinationDomain,
		Nonce:             nonce,
		Sender:            messageSender,
		Recipient:         recipient,
		DestinationCaller: destinationCaller,
		MessageBody:       messageBody,
	}

	messageBytes, err := message.Bytes()
	if err != nil {
		return err
	}
	event := types.MessageSent{
		Message: messageBytes,
	}
	err = ctx.EventManager().EmitTypedEvent(&event)

	return err
}
