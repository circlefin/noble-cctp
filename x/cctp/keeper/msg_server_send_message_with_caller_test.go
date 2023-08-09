package keeper_test

import (
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
 * Happy path
 * Invalid destination caller
 * Sending and receiving messages is paused
 * Recipient is empty
 * Message body is too long
 */
func TestSendMessageWithCallerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessageWithCaller{
		From:              "12345678901234567890123456789012",
		DestinationDomain: 3,
		Recipient:         []byte("12345678901234567890123456789012"),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	resp, err := server.SendMessageWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.Nil(t, err)
	require.Equal(t, nonce.Nonce, resp.Nonce)

	nextNonce, found := testkeeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t, nonce.Nonce+1, nextNonce.Nonce)
}

func TestSendMessageWithCallerInvalidDestinationCaller(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessageWithCaller{
		From:              "12345678901234567890123456789012",
		DestinationDomain: 3,
		Recipient:         []byte("12345678901234567890123456789012"),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
	}

	_, err := server.SendMessageWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
	require.Contains(t, err.Error(), "destination caller must be nonzero")
}

func TestSendMessageWithCallerSendingAndReceivingMessagesPaused(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: true}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	msg := types.MsgSendMessageWithCaller{
		From:              "12345678901234567890123456789012",
		DestinationDomain: 3,
		Recipient:         []byte("12345678901234567890123456789012"),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.SendMessageWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
	require.Contains(t, err.Error(), "sending and receiving messages is paused")
}

func TestSendMessageWithCallerRecipientEmpty(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessageWithCaller{
		From:              "anything",
		DestinationDomain: 3,
		Recipient:         make([]byte, types.MintRecipientLen),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.SendMessageWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
	require.Contains(t, err.Error(), "recipient must not be nonzero")
}

func TestSendMessageWithCallerMessageBodyTooLong(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	max := types.MaxMessageBodySize{Amount: 5}
	testkeeper.SetMaxMessageBodySize(ctx, max)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessageWithCaller{
		From:              "anything",
		DestinationDomain: 3,
		Recipient:         []byte("12345678901234567890123456789012"),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.SendMessageWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
	require.Contains(t, err.Error(), "message body exceeds max size")

}
