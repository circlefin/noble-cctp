package keeper

import (
	"bytes"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// DecodeMessage parses a byte array into a Message struct
// https://developers.circle.com/stablecoin/docs/cctp-technical-reference#message
func DecodeMessage(msg []byte) types.Message {
	message := types.Message{
		Version:           binary.BigEndian.Uint32(msg[types.VersionIndex:types.SourceDomainIndex]),
		SourceDomain:      binary.BigEndian.Uint32(msg[types.SourceDomainIndex:types.DestinationDomainIndex]),
		DestinationDomain: binary.BigEndian.Uint32(msg[types.DestinationDomainIndex:types.NonceIndex]),
		Nonce:             binary.BigEndian.Uint64(msg[types.NonceIndex:types.SenderIndex]),
		Sender:            msg[types.SenderIndex:types.RecipientIndex],
		Recipient:         msg[types.RecipientIndex:types.DestinationCallerIndex],
		DestinationCaller: msg[types.DestinationCallerIndex:types.MessageBodyIndex],
		MessageBody:       msg[types.MessageBodyIndex:],
	}

	return message
}

// DecodeBurnMessage parses a byte array into a BurnMessage struct
// https://developers.circle.com/stablecoin/docs/cctp-technical-reference#burnmessage
func DecodeBurnMessage(msg []byte) *types.BurnMessage {
	message := types.BurnMessage{
		Version:       binary.BigEndian.Uint32(msg[types.BurnMsgVersionIndex:types.BurnTokenIndex]),
		BurnToken:     msg[types.BurnTokenIndex:types.MintRecipientIndex],
		MintRecipient: msg[types.MintRecipientIndex:types.AmountIndex],
		Amount:        math.NewIntFromBigInt(new(big.Int).SetBytes(msg[types.AmountIndex:types.MsgSenderIndex])),
		MessageSender: msg[types.MsgSenderIndex:types.BurnMessageLen],
	}

	return &message
}

// EncodeMessage parses a Message struct into a byte array
// sender, recipient, destination caller must be 32 bytes
func EncodeMessage(msg types.Message) ([]byte, error) {
	if len(msg.Sender) != types.AddressBytesLen ||
		len(msg.Recipient) != types.AddressBytesLen ||
		len(msg.DestinationCaller) != types.AddressBytesLen {
		return nil, sdkerrors.Wrap(types.ErrParsingMessage, "sender, recipient, destination caller must be 32 bytes")
	}
	result := make([]byte, types.MessageBodyIndex+len(msg.MessageBody))

	versionBytes := make([]byte, types.VersionLen)
	binary.BigEndian.PutUint32(versionBytes, msg.Version)

	sourceDomainBytes := make([]byte, types.DomainBytesLen)
	binary.BigEndian.PutUint32(sourceDomainBytes, msg.SourceDomain)

	destinationDomainBytes := make([]byte, types.DomainBytesLen)
	binary.BigEndian.PutUint32(destinationDomainBytes, msg.DestinationDomain)

	nonceBytes := make([]byte, types.NonceBytesLen)
	binary.BigEndian.PutUint64(nonceBytes, msg.Nonce)

	copyBytes(types.VersionIndex, types.SourceDomainIndex, versionBytes, &result)
	copyBytes(types.SourceDomainIndex, types.DestinationDomainIndex, sourceDomainBytes, &result)
	copyBytes(types.DestinationDomainIndex, types.NonceIndex, destinationDomainBytes, &result)
	copyBytes(types.NonceIndex, types.SenderIndex, nonceBytes, &result)
	copyBytes(types.SenderIndex, types.RecipientIndex, msg.Sender, &result)
	copyBytes(types.RecipientIndex, types.DestinationCallerIndex, msg.Recipient, &result)
	copyBytes(types.DestinationCallerIndex, types.MessageBodyIndex, msg.DestinationCaller, &result)
	copyBytes(types.MessageBodyIndex, types.MessageBodyIndex+len(msg.MessageBody), msg.MessageBody, &result)

	return result, nil
}

// EncodeBurnMessage parses a BurnMessage struct into a byte array
// burn token, mint recipient, and message sender must be 32 bytes
func EncodeBurnMessage(msg types.BurnMessage) ([]byte, error) {
	if len(msg.BurnToken) != types.BurnTokenLen ||
		len(msg.MintRecipient) != types.MintRecipientLen ||
		len(msg.MessageSender) != types.AddressBytesLen {
		return nil, sdkerrors.Wrap(types.ErrParsingBurnMessage, "burn token, mint recipient, message sender must be 32 bytes")
	}

	result := make([]byte, types.BurnMessageLen)

	versionBytes := make([]byte, types.VersionLen)
	binary.BigEndian.PutUint32(versionBytes, msg.Version)

	amountBytes := make([]byte, types.AmountLen)
	msg.Amount.BigInt().FillBytes(amountBytes)

	copyBytes(types.BurnMsgVersionIndex, types.BurnTokenIndex, versionBytes, &result)
	copyBytes(types.BurnTokenIndex, types.MintRecipientIndex, msg.BurnToken, &result)
	copyBytes(types.MintRecipientIndex, types.AmountIndex, msg.MintRecipient, &result)
	copyBytes(types.AmountIndex, types.MsgSenderIndex, amountBytes, &result)
	copyBytes(types.MsgSenderIndex, types.BurnMessageLen, msg.MessageSender, &result)

	return result, nil
}

/**
 * Copies the contents of the copyFrom byte slice into the copyInto byte slice,
 * starting at the position start and ending at the position end in the copyInto slice
 */
func copyBytes(start int, end int, copyFrom []byte, copyInto *[]byte) {
	for i := 0; i < end-start; i++ {
		(*copyInto)[i+start] = copyFrom[i]
	}
}

/**
 * VerifyAttestationSignatures return true if a message was signed by enough private keys
 * @param message the MessageSent message bytes
 * @param attestation a concatenated list of message signatures
 * @param publicKeys a list of hex encoded ECDSA uncompressed public keys used to verify message signatures
 * @param signatureThreshold the minimum amount of signatures in an attestation to consider it valid
 */
func VerifyAttestationSignatures(
	message []byte,
	attestation []byte,
	publicKeys []types.Attester,
	signatureThreshold uint32) (bool, error) {

	/*
	* Rules for valid attestation:
	* 1. length of `_attestation` == 65 (signature length) * signatureThreshold
	* 2. addresses recovered from attestation must be in increasing order.
	* 	For example, if signature A is signed by address 0x1..., and signature B
	* 		is signed by address 0x2..., attestation must be passed as AB.
	* 3. no duplicate signers
	* 4. all signers must be enabled attesters
	 */

	if uint32(len(attestation)) != types.SignatureLength*signatureThreshold {
		return false, sdkerrors.Wrap(types.ErrSignatureVerification, "invalid attestation length")
	}

	if signatureThreshold == 0 {
		return false, sdkerrors.Wrap(types.ErrSignatureVerification, "signature verification threshold cannot be 0")
	}

	// public keys cannot be empty, so the recovered key should be bigger than latestECDSA
	var latestECDSA ecdsa.PublicKey

	digest := crypto.Keccak256(message)

	for i := uint32(0); i < signatureThreshold; i++ {
		signature := attestation[i*types.SignatureLength : (i*types.SignatureLength)+types.SignatureLength]

		recoveredKey, err := crypto.Ecrecover(digest, signature)
		if err != nil {
			return false, sdkerrors.Wrap(types.ErrSignatureVerification, "failed to recover public key")
		}

		// Signatures must be in increasing order of address, and may not duplicate signatures from same address
		recoveredECSDA := ecdsa.PublicKey{
			X: new(big.Int).SetBytes(recoveredKey[1:33]),
			Y: new(big.Int).SetBytes(recoveredKey[33:]),
		}

		if latestECDSA.X != nil && latestECDSA.Y != nil && bytes.Compare(
			crypto.PubkeyToAddress(latestECDSA).Bytes(),
			crypto.PubkeyToAddress(recoveredECSDA).Bytes()) > -1 {
			return false, sdkerrors.Wrap(types.ErrSignatureVerification, "invalid signature order or dupe")
		}

		// check that recovered key is a valid
		contains := false
		for _, key := range publicKeys {
			hexBz, err := hex.DecodeString(key.Attester)
			if err != nil {
				return false, sdkerrors.Wrap(types.ErrSignatureVerification, "invalid signature: not attester")
			}
			if bytes.Equal(hexBz, recoveredKey) {
				contains = true
				break
			}
		}

		if !contains {
			return false, sdkerrors.Wrap(types.ErrSignatureVerification, "Invalid signature: not an attester")
		}

		latestECDSA.X, latestECDSA.Y = recoveredECSDA.X, recoveredECSDA.Y
	}
	return true, nil
}
