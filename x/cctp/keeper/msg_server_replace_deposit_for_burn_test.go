package keeper_test

import (
	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"testing"
)

/**
 * Happy path
 * Fails when paused
 * Outer message too short
 * Burn message invalid length
 * Invalid sender
 * Invalid new mint recipient
 */

func TestReplaceDepositForBurnHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.BurningAndMintingPaused{Paused: false}
	testkeeper.SetBurningAndMintingPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := burnMessage.Bytes()
	require.NoError(t, err)

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
	originalMessageBytes, err := originalMessage.Bytes()
	require.NoError(t, err)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceDepositForBurn{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     []byte("new mint recipient90123456789012"),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.Nil(t, err)
}

func TestReplaceDepositForBurnFailsWhenPaused(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.BurningAndMintingPaused{Paused: false}
	testkeeper.SetBurningAndMintingPaused(ctx, paused)

	_, err := server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &types.MsgReplaceDepositForBurn{})
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "sending and receiving messages are paused")
}

func TestReplaceDepositForBurnOuterMessageTooShort(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.BurningAndMintingPaused{Paused: false}
	testkeeper.SetBurningAndMintingPaused(ctx, paused)

	_, err := server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &types.MsgReplaceDepositForBurn{})
	require.ErrorIs(t, types.ErrParsingMessage, err)
	require.Contains(t, err.Error(), "cctp message must be at least 116 bytes, got 0: error while parsing message into bytes")
}

func TestReplaceDepositForBurnBurnMessageInvalidLength(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.BurningAndMintingPaused{Paused: false}
	testkeeper.SetBurningAndMintingPaused(ctx, paused)

	burnMessageBytes := make([]byte, types.BurnMessageLen+1)

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
	originalMessageBytes, err := originalMessage.Bytes()
	require.NoError(t, err)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceDepositForBurn{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     []byte("new mint recipient90123456789012"),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "burn message body is not the correct length")
}

func TestReplaceDepositForBurnInvalidSender(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.BurningAndMintingPaused{Paused: false}
	testkeeper.SetBurningAndMintingPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := burnMessage.Bytes()
	require.NoError(t, err)

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
	originalMessageBytes, err := originalMessage.Bytes()
	require.NoError(t, err)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceDepositForBurn{
		From:                 "not the original sender",
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     []byte("new mint recipient90123456789012"),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "invalid sender for message")
}

func TestReplaceDepositForBurnInvalidNewMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: make([]byte, 32),
	}

	burnMessageBytes, err := burnMessage.Bytes()
	require.NoError(t, err)

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
	originalMessageBytes, err := originalMessage.Bytes()
	require.NoError(t, err)

	// generate attestation, set attesters, signature threshold
	signatureThreshold := uint32(2)
	privKeys := generateNPrivateKeys(int(signatureThreshold))
	attesters := getAttestersFromPrivateKeys(privKeys)
	originalAttestation := generateAttestation(originalMessageBytes, privKeys)
	for _, attester := range attesters {
		testkeeper.SetAttester(ctx, attester)
	}
	testkeeper.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: signatureThreshold})

	msg := types.MsgReplaceDepositForBurn{
		From:                 string(originalMessage.Sender),
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     make([]byte, types.MintRecipientLen),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}
