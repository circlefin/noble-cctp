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
package keeper

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"

	sdkerrors "cosmossdk.io/errors"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	xIndex = 1
	yIndex = 33
)

/*
* Rules for valid attestation:
* 1. length of `_attestation` == 65 (signature length) * signatureThreshold
* 2. addresses recovered from attestation must be in increasing order.
* 	For example, if signature A is signed by address 0x1..., and signature B
* 		is signed by address 0x2..., attestation must be passed as AB.
* 3. no duplicate signers
* 4. all signers must be enabled attesters
 */
func VerifyAttestationSignatures(
	message []byte,
	attestation []byte,
	publicKeys []types.Attester,
	signatureThreshold uint32,
) error {
	if uint32(len(attestation)) != types.SignatureLength*signatureThreshold {
		return sdkerrors.Wrap(types.ErrSignatureVerification, "invalid attestation length")
	}

	if signatureThreshold == 0 {
		return sdkerrors.Wrap(types.ErrSignatureVerification, "signature verification threshold cannot be 0")
	}

	// public keys cannot be empty, so the recovered key should be bigger than latestECDSA
	var latestECDSA ecdsa.PublicKey

	digest := crypto.Keccak256(message)

	for i := uint32(0); i < signatureThreshold; i++ {
		signature := attestation[i*types.SignatureLength : (i*types.SignatureLength)+types.SignatureLength]

		// The go-ethereum library assumes that the v-value of a secp256k1
		// signature is either 0 or 1. However, in legacy Bitcoin signers, this
		// value is either 27 or 28. So we need to normalise in order to verify.
		if signature[len(signature)-1] == 27 || signature[len(signature)-1] == 28 {
			signature[len(signature)-1] -= 27
		}

		recoveredKey, err := crypto.Ecrecover(digest, signature)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrSignatureVerification, "failed to recover public key: %s", err)
		}

		// Signatures must be in increasing order of address, and may not duplicate signatures from same address
		recoveredECSDA := ecdsa.PublicKey{
			X: new(big.Int).SetBytes(recoveredKey[xIndex:yIndex]),
			Y: new(big.Int).SetBytes(recoveredKey[yIndex:]),
		}

		if latestECDSA.X != nil && latestECDSA.Y != nil && bytes.Compare(
			crypto.PubkeyToAddress(latestECDSA).Bytes(),
			crypto.PubkeyToAddress(recoveredECSDA).Bytes()) > -1 {
			return sdkerrors.Wrap(types.ErrSignatureVerification, "invalid signature order or dupe")
		}

		// check that recovered key is a valid attester
		contains := false
		for _, key := range publicKeys {
			hexBz := common.FromHex(key.Attester)
			if bytes.Equal(hexBz, recoveredKey) {
				contains = true
				break
			}
		}

		if !contains {
			return sdkerrors.Wrap(types.ErrSignatureVerification, "Invalid signature: not an attester")
		}

		latestECDSA.X, latestECDSA.Y = recoveredECSDA.X, recoveredECSDA.Y
	}
	return nil
}
