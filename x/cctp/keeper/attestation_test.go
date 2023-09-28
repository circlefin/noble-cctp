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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"log"
	"sort"
	"testing"

	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestVerifyAttestationSignaturesHappyPath(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(66)
	attesters := getAttestersFromPrivateKeys(privKeys)
	attestation := generateAttestation(message, privKeys)

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 66)
	require.NoError(t, err)
}

func TestVerifyAttestationSignaturesWithSmallerThresholdThanAttesterCount(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(66)
	attestation := generateAttestation(message, privKeys)

	// generate more attesters that won't be used
	morePrivKeys := generateNPrivateKeys(120)
	attesters := append(getAttestersFromPrivateKeys(privKeys), getAttestersFromPrivateKeys(morePrivKeys)...)

	// signature threshold < attesters
	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 66)
	require.NoError(t, err)
}

func TestVerifyAttestationSignaturesInvalidAttestationLength(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(66)
	attesters := getAttestersFromPrivateKeys(privKeys)
	attestation := []byte("an attestation that i")

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 66)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "invalid attestation length")
}

func TestVerifyAttestationSignaturesSignatureThresholdIsZero(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(2)
	attesters := getAttestersFromPrivateKeys(privKeys)
	var attestation []byte

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 0)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "signature verification threshold cannot be 0")
}

func TestVerifyAttestationSignaturesFailedToRecoverPublicKey(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(2)
	attesters := getAttestersFromPrivateKeys(privKeys)
	differentPrivKeys := generateNPrivateKeys(2)
	attestation := generateAttestation(message, differentPrivKeys)
	attestation[64] = 5 // Invalid recovery ID

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 2)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "failed to recover public key")
}

func TestVerifyAttestationSignaturesInvalidSignatureOrder(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(20000) // high number to increase odds of invalid sort order
	attesters := getAttestersFromPrivateKeys(privKeys)
	attestation := generateAttestationWithInvalidSignatureOrder(message, privKeys)

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 20000)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "invalid signature order or dupe")
}

func TestVerifyAttestationSignaturesDupe(t *testing.T) {
	message := []byte("Execute order")
	privKeys := generateNPrivateKeys(2)
	attesters := getAttestersFromPrivateKeys(privKeys)
	attestation := generateAttestationWithDupe(message, privKeys)

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, 3)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "invalid signature order or dupe")
}

func TestVerifyAttestationSignaturesLegacySignature(t *testing.T) {
	message := []byte("")

	keys := generateNPrivateKeys(2)
	attesters := getAttestersFromPrivateKeys(keys)
	attestation := generateAttestation(message, keys)

	// Because we use updated signers, we need to mock a legacy signature by
	// manually changing the v-value.
	for i := 0; i < len(keys); i++ {
		signature := attestation[i*types.SignatureLength : (i*types.SignatureLength)+types.SignatureLength]
		signature[len(signature)-1] += 27
	}

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters, uint32(len(keys)))
	require.NoError(t, err)
}

func TestVerifyAttestationSignaturesInvalidAttester(t *testing.T) {
	message := []byte("")

	keys := generateNPrivateKeys(3)
	attesters := getAttestersFromPrivateKeys(keys)
	attestation := generateAttestation(message, keys)

	attestation = append(attestation[0:types.SignatureLength], attestation[2*types.SignatureLength:]...)

	err := keeper.VerifyAttestationSignatures(message, attestation, attesters[:1], 2)
	require.ErrorIs(t, types.ErrSignatureVerification, err)
	require.Contains(t, err.Error(), "Invalid signature: not an attester")
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
		publicKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y) //nolint:staticcheck

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

func generateAttestationWithDupe(message []byte, privKeys []*ecdsa.PrivateKey) (attestation []byte) {
	type Attestation struct {
		pubKey      ecdsa.PublicKey
		attestation []byte // 65 byte
	}
	attestationList := make([]Attestation, len(privKeys)+1)

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

	attestationList[len(privKeys)] = attestationList[len(privKeys)-1]

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
