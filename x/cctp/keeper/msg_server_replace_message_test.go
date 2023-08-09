package keeper_test

import (
	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
 * Happy path
 * Fails when paused
 * Signature threshold not found
 * Signature verification failed
 * Message body too short
 * Invalid sender
 * Message not originally sent from this domain
 */

func TestReplaceMessageHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := keeper.EncodeBurnMessage(burnMessage)
	require.Nil(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       burnMessageBytes,
	}
	originalMessageBytes, err := keeper.EncodeMessage(originalMessage)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceMessage{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewMessageBody:       []byte("123"),
		NewDestinationCaller: []byte("new destination caller3456789012"),
	}

	_, err = server.ReplaceMessage(sdk.WrapSDKContext(ctx), &msg)
	require.Nil(t, err)
}

func TestReplaceMessageFailsWhenPaused(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: true}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	_, err := server.ReplaceMessage(sdk.WrapSDKContext(ctx), &types.MsgReplaceMessage{})
	require.ErrorIs(t, types.ErrReplaceMessage, err)
	require.Contains(t, err.Error(), "sending and receiving messages are paused")
}

func TestReplaceMessageSignatureThresholdNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := keeper.EncodeBurnMessage(burnMessage)
	require.Nil(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       burnMessageBytes,
	}
	originalMessageBytes, err := keeper.EncodeMessage(originalMessage)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}

	msg := types.MsgReplaceMessage{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewMessageBody:       []byte("123"),
		NewDestinationCaller: []byte("new destination caller3456789012"),
	}

	_, err = server.ReplaceMessage(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrReplaceMessage, err)
	require.Contains(t, err.Error(), "signature threshold not found")
}

func TestReplaceMessageSignatureVerificationFailed(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := keeper.EncodeBurnMessage(burnMessage)
	require.Nil(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       burnMessageBytes,
	}
	originalMessageBytes, err := keeper.EncodeMessage(originalMessage)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	// corrupt attestation
	originalAttestation[10] = 1

	msg := types.MsgReplaceMessage{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewMessageBody:       []byte("123"),
		NewDestinationCaller: []byte("new destination caller3456789012"),
	}

	_, err = server.ReplaceMessage(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "unable to verify signatures")
}

func TestReplaceMessageMessageBodyTooShort(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := keeper.EncodeBurnMessage(burnMessage)
	require.Nil(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       burnMessageBytes,
	}
	originalMessageBytes, err := keeper.EncodeMessage(originalMessage)
	// make it too small
	originalMessageBytes = originalMessageBytes[0:115]

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceMessage{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewMessageBody:       []byte("123"),
		NewDestinationCaller: []byte("new destination caller3456789012"),
	}

	_, err = server.ReplaceMessage(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrReplaceMessage, err)
	require.Contains(t, err.Error(), "invalid message: too short")
}

func TestReplaceMessageInvalidSender(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := keeper.EncodeBurnMessage(burnMessage)
	require.Nil(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       burnMessageBytes,
	}
	originalMessageBytes, err := keeper.EncodeMessage(originalMessage)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceMessage{
		From:                 "not the original sender",
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewMessageBody:       []byte("123"),
		NewDestinationCaller: []byte("new destination caller3456789012"),
	}

	_, err = server.ReplaceMessage(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrReplaceMessage, err)
	require.Contains(t, err.Error(), "sender not permitted to use nonce")
}

func TestReplaceMessageMessageNotOriginallySentFromThisDomain(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := keeper.EncodeBurnMessage(burnMessage)
	require.Nil(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      8, // not Noble's domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       burnMessageBytes,
	}
	originalMessageBytes, err := keeper.EncodeMessage(originalMessage)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceMessage{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewMessageBody:       []byte("123"),
		NewDestinationCaller: []byte("new destination caller3456789012"),
	}

	_, err = server.ReplaceMessage(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrReplaceMessage, err)
	require.Contains(t, err.Error(), "message not originally sent from this domain")
}
