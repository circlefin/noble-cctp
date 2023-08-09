package keeper_test

import (
	"bytes"
	"cosmossdk.io/math"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"log"
	"sort"
	"testing"
)

// Message -> bytes -> Message -> bytes
func TestParseMessageHappyPath(t *testing.T) {
	message := types.Message{
		Version:           1,
		SourceDomain:      2,
		DestinationDomain: 3,
		Nonce:             4,
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
		MessageBody:       []byte("message body"),
	}
	messageBytes, err := keeper.EncodeMessage(message)
	require.Nil(t, err)
	parsedMessage := keeper.DecodeMessage(messageBytes)
	require.Equal(t, message, parsedMessage)
	parsedMessageBytes, err := keeper.EncodeMessage(parsedMessage)
	require.Nil(t, err)
	require.Equal(t, messageBytes, parsedMessageBytes)
}

func TestParseIntoMessageWithInvalidInput(t *testing.T) {
	message := types.Message{
		Sender:            []byte("too short"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("destination caller90123456789012"),
	}
	_, err := keeper.EncodeMessage(message)
	require.ErrorIs(t, types.ErrParsingMessage, err)

	message = types.Message{
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("too short"),
		DestinationCaller: []byte("destination caller90123456789012"),
	}
	_, err = keeper.EncodeMessage(message)
	require.ErrorIs(t, types.ErrParsingMessage, err)

	message = types.Message{
		Sender:            []byte("sender78901234567890123456789012"),
		Recipient:         []byte("recipient01234567890123456789012"),
		DestinationCaller: []byte("too short"),
	}
	_, err = keeper.EncodeMessage(message)
	require.ErrorIs(t, types.ErrParsingMessage, err)
}

// BurnMessage -> bytes -> BurnMessage -> bytes
func TestParseIntoBurnMessageHappyPath(t *testing.T) {
	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     crypto.Keccak256([]byte("usdc")),
		MintRecipient: []byte("recipient01234567890123456789012"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("message-sender567890123456789012"),
	}
	burnMessageBytes, err := keeper.EncodeBurnMessage(burnMessage)
	require.Nil(t, err)
	parsedBurnMessage := keeper.DecodeBurnMessage(burnMessageBytes)

	require.Equal(t, burnMessage.Version, parsedBurnMessage.Version)
	require.Equal(t, burnMessage.BurnToken, parsedBurnMessage.BurnToken)
	require.Equal(t, burnMessage.MintRecipient, parsedBurnMessage.MintRecipient)
	require.Equal(t, burnMessage.Amount, parsedBurnMessage.Amount)
	require.Equal(t, burnMessage.MessageSender, parsedBurnMessage.MessageSender)

	parsedBurnMessageBytes, err := keeper.EncodeBurnMessage(*parsedBurnMessage)
	require.Nil(t, err)
	require.Equal(t, burnMessageBytes, parsedBurnMessageBytes)
}

func TestParseIntoBurnMessageWithInvalidInput(t *testing.T) {
	burnMessage := types.BurnMessage{
		Version:       1,
		BurnToken:     []byte("too short"),
		MintRecipient: []byte("recipient01234567890123456789012"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("message-sender567890123456789012"),
	}
	_, err := keeper.EncodeBurnMessage(burnMessage)
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)

	burnMessage = types.BurnMessage{
		Version:       1,
		BurnToken:     crypto.Keccak256([]byte("usdc")),
		MintRecipient: []byte("too short"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("message-sender567890123456789012"),
	}
	_, err = keeper.EncodeBurnMessage(burnMessage)
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)

	burnMessage = types.BurnMessage{
		Version:       1,
		BurnToken:     crypto.Keccak256([]byte("usdc")),
		MintRecipient: []byte("recipient01234567890123456789012"),
		Amount:        math.NewInt(345678),
		MessageSender: []byte("too short"),
	}
	_, err = keeper.EncodeBurnMessage(burnMessage)
	require.ErrorIs(t, types.ErrParsingBurnMessage, err)
}

func TestVerifyAttestationSignaturesHappyPath(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(66)
	attesters := getAttestersFromPrivateKeys(privKeys)
	attestation := generateAttestation(message, privKeys)

	verified, err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 66)
	require.Nil(t, err)
	require.True(t, verified)
}

func TestVerifyAttestationSignaturesWithSmallerThresholdThanAttesterCount(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(66)
	attestation := generateAttestation(message, privKeys)

	// generate more attesters that won't be used
	morePrivKeys := generateNPrivateKeys(120)
	attesters := append(getAttestersFromPrivateKeys(privKeys), getAttestersFromPrivateKeys(morePrivKeys)...)

	// signature threshold < attesters
	verified, err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 66)
	require.Nil(t, err)
	require.True(t, verified)
}

func TestVerifyAttestationSignaturesInvalidAttestationLength(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(66)
	attesters := getAttestersFromPrivateKeys(privKeys)
	attestation := []byte("an attestation that i")

	verified, err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 66)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "invalid attestation length")
	require.False(t, verified)
}

func TestVerifyAttestationSignaturesSignatureThresholdIsZero(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(2)
	attesters := getAttestersFromPrivateKeys(privKeys)
	var attestation []byte

	verified, err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 0)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "signature verification threshold cannot be 0")
	require.False(t, verified)
}

func TestVerifyAttestationSignaturesFailedToRecoverPublicKey(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(2)
	attesters := getAttestersFromPrivateKeys(privKeys)
	differentPrivKeys := generateNPrivateKeys(2)
	attestation := generateAttestation(message, differentPrivKeys)
	attestation[64] = 5 // Invalid recovery ID

	verified, err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 2)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "failed to recover public key")
	require.False(t, verified)
}

func TestVerifyAttestationSignaturesInvalidSignatureOrder(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(20000) // high number to increase odds of invalid sort order
	attesters := getAttestersFromPrivateKeys(privKeys)
	attestation := generateAttestationWithInvalidSignatureOrder(message, privKeys)

	verified, err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 20000)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "invalid signature order or dupe")
	require.False(t, verified)
}

func generateNPrivateKeys(n int) []*ecdsa.PrivateKey {
	result := make([]*ecdsa.PrivateKey, n)
	for i := 0; i < n; i++ {
		result[i], _ = crypto.GenerateKey()
	}
	return result
}

func getAttestersFromPrivateKeys(privkeys []*ecdsa.PrivateKey) []types.Attester {
	result := make([]types.Attester, len(privkeys))
	for i, privkey := range privkeys {
		// Get the public key
		publicKey := privkey.PublicKey

		// Marshal the public key into bytes
		publicKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)

		result[i] = types.Attester{Attester: hex.EncodeToString(publicKeyBytes)}
	}
	return result
}

func generateAttestation(message []byte, privKeys []*ecdsa.PrivateKey) (attestation []byte) {
	type Attestation struct {
		pubKey      ecdsa.PublicKey
		attestation []byte // 65 byte
	}
	attestationList := make([]Attestation, len(privKeys))

	for i, privateKey := range privKeys {
		// Sign the message with the private key
		sig, err := crypto.Sign(crypto.Keccak256Hash(message).Bytes(), privateKey)
		if err != nil {
			log.Fatalf("Failed to sign message: %v", err)
		}
		attestationList[i] = Attestation{
			pubKey:      privateKey.PublicKey,
			attestation: sig,
		}
	}

	sort.Slice(attestationList, func(i, j int) bool {
		return bytes.Compare(
			crypto.PubkeyToAddress(attestationList[i].pubKey).Bytes(),
			crypto.PubkeyToAddress(attestationList[j].pubKey).Bytes(),
		) < 0
	})

	var result []byte
	for _, att := range attestationList {
		result = append(result, att.attestation...)
	}

	return result
}

func generateAttestationWithInvalidSignatureOrder(message []byte, privKeys []*ecdsa.PrivateKey) (attestation []byte) {
	type Attestation struct {
		pubKey      ecdsa.PublicKey
		attestation []byte // 65 byte
	}
	attestationList := make([]Attestation, len(privKeys))

	for i, privateKey := range privKeys {
		// Sign the message with the private key
		sig, err := crypto.Sign(crypto.Keccak256Hash(message).Bytes(), privateKey)
		if err != nil {
			log.Fatalf("Failed to sign message: %v", err)
		}
		attestationList[i] = Attestation{
			pubKey:      privateKey.PublicKey,
			attestation: sig,
		}
	}

	var result []byte
	for _, att := range attestationList {
		result = append(result, att.attestation...)
	}

	return result
}
