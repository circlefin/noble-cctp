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
package cli

import (
	"errors"

	"github.com/cosmos/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
)

// parseAddress parses an encoded address into a 32 length byte array.
// Currently supported encodings: base58, hex.
func parseAddress(address string) ([]byte, error) {
	if address[:2] == "0x" {
		bz := common.FromHex(address)
		return leftPadBytes(bz)
	}

	bz := base58.Decode(address)
	return leftPadBytes(bz)
}

// leftPadBytes left pads a byte array to be length 32.
func leftPadBytes(bz []byte) ([]byte, error) {
	res := make([]byte, 32)

	if len(bz) > 32 {
		return nil, errors.New("decoded bytes too big")
	}

	copy(res[32-len(bz):], bz)
	return res, nil
}
