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
	"testing"

	"github.com/circlefin/noble-cctp/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgSendMessageWithCaller_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSendMessageWithCaller
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgSendMessageWithCaller{
				From:              "invalid_address",
				DestinationDomain: 123,
				Recipient:         []byte{2, 3, 4},
				MessageBody:       []byte{2, 3, 4},
				DestinationCaller: []byte{2, 3, 4},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgSendMessageWithCaller{
				From:              sample.AccAddress(),
				DestinationDomain: 123,
				Recipient:         []byte{2, 3, 4},
				MessageBody:       []byte{2, 3, 4},
				DestinationCaller: []byte{2, 3, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
