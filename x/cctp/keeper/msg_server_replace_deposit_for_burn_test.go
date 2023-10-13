/*
 * Copyright (c) 2023, © Circle Internet Financial, LTD.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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

	// we encode the message sender when sending messages, so we must use an encoded message in the original message
	sender := sample.AccAddress()
	senderEncoded := make([]byte, 32)
	copy(senderEncoded[12:], sdk.MustAccAddressFromBech32(sender))

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: senderEncoded,
	}

	burnMessageBytes, err := burnMessage.Bytes()
	require.NoError(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            types.PaddedModuleAddress,
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
		From:                 sender,
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

	paused := types.BurningAndMintingPaused{Paused: true}
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
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "burning and minting are paused")
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
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)
	require.Contains(t, err.Error(), "burn message must be 132 bytes")
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

	// we encode the message sender when sending messages, so we must use an encoded message in the original message
	sender := sample.AccAddress()
	senderEncoded := make([]byte, 32)
	copy(senderEncoded[12:], sdk.MustAccAddressFromBech32(sender))

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            senderEncoded, // different sender than the replaceMessage sender
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
		From:                 sample.AccAddress(), // different sender
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     []byte("new mint recipient90123456789012"),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "invalid sender for message")
}

func TestReplaceDepositForBurnEmptyNewMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	// we encode the message sender when sending messages, so we must use an encoded message in the original message
	sender := sample.AccAddress()
	senderEncoded := make([]byte, 32)
	copy(senderEncoded[12:], sdk.MustAccAddressFromBech32(sender))

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: senderEncoded,
	}

	burnMessageBytes, err := burnMessage.Bytes()
	require.NoError(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            types.PaddedModuleAddress,
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
		From:                 sender,
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     make([]byte, types.MintRecipientLen),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}

func TestReplaceDepositForBurnInvalidNewMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	// we encode the message sender when sending messages, so we must use an encoded message in the original message
	sender := sample.AccAddress()
	senderEncoded := make([]byte, 32)
	copy(senderEncoded[12:], sdk.MustAccAddressFromBech32(sender))

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: senderEncoded,
	}

	burnMessageBytes, err := burnMessage.Bytes()
	require.NoError(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      4, // Noble domain id
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            types.PaddedModuleAddress,
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
		From:                 sender,
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     []byte("INVALID RECIPIENT"),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)
	require.Contains(t, err.Error(), "error parsing burn message")
}

func TestReplaceDepositForBurnIncorrectSourceID(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.BurningAndMintingPaused{Paused: false}
	testkeeper.SetBurningAndMintingPaused(ctx, paused)

	// we encode the message sender when sending messages, so we must use an encoded message in the original message
	sender := sample.AccAddress()
	senderEncoded := make([]byte, 32)
	copy(senderEncoded[12:], sdk.MustAccAddressFromBech32(sender))

	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     make([]byte, 32),
		MintRecipient: make([]byte, 32),
		Amount:        math.NewInt(123456),
		MessageSender: senderEncoded,
	}

	burnMessageBytes, err := burnMessage.Bytes()
	require.NoError(t, err)

	originalMessage := types.Message{
		Version:           1,
		SourceDomain:      9000, // NOT THE NOBLE DOMAIN ID
		DestinationDomain: 3,
		Nonce:             2,
		Sender:            types.PaddedModuleAddress,
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
		From:                 sender,
		OriginalMessage:      originalMessageBytes,
		OriginalAttestation:  originalAttestation,
		NewDestinationCaller: []byte("new destination caller3456789012"),
		NewMintRecipient:     []byte("new mint recipient90123456789012"),
	}

	_, err = server.ReplaceDepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "message not originally sent from this domain")
}
