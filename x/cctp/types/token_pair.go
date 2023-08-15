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
package types

import (
	"encoding/hex"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RemoteTokenPadded returns the remote token as a byte array, padded to 32 bytes
func RemoteTokenPadded(remoteTokenHex string) ([]byte, error) {
	remoteToken, err := hex.DecodeString(strings.TrimPrefix(remoteTokenHex, "0x"))
	if err != nil {
		return nil, sdkerrors.Wrap(ErrInvalidRemoteToken, "must be hex string")
	}

	if len(remoteToken) > BurnTokenLen {
		return nil, sdkerrors.Wrapf(ErrInvalidRemoteToken, "must be less than %d bytes", BurnTokenLen)
	}

	remoteTokenPadded := make([]byte, BurnTokenLen)
	for i := 0; i < BurnTokenLen-len(remoteToken); i++ {
		remoteTokenPadded[i] = 0
	}
	copy(remoteTokenPadded[BurnTokenLen-len(remoteToken):], remoteToken)

	return remoteTokenPadded, nil
}
