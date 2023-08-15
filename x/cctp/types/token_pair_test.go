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
package types_test

import (
	"testing"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/stretchr/testify/require"
)

func TestRemoteTokenPadded(t *testing.T) {
	type tc struct {
		name           string
		remoteTokenHex string
		expected       []byte
		err            error
	}

	tcs := []tc{
		{
			name:           "happy path",
			remoteTokenHex: "0xabcd",
			expected:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xab, 0xcd},
		},
		{
			name:           "happy path no 0x",
			remoteTokenHex: "abcd",
			expected:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xab, 0xcd},
		},
		{
			name:           "overflow",
			remoteTokenHex: "0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			err:            types.ErrInvalidRemoteToken,
		},
		{
			name:           "invalid hex",
			remoteTokenHex: "invalid",
			err:            types.ErrInvalidRemoteToken,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			bz, err := types.RemoteTokenPadded(tc.remoteTokenHex)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, bz)
			}
		})
	}
}
