// Copyright 2024 Circle Internet Group, Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/binary"
	"math/big"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
)

// Parse parses a byte array into a BurnMessage struct
// https://developers.circle.com/stablecoin/docs/cctp-technical-reference#burnmessage
func (msg *BurnMessage) Parse(bz []byte) (*BurnMessage, error) {
	if len(bz) != BurnMessageLen {
		return nil, errors.Wrapf(ErrParsingBurnMessage, "burn message must be %d bytes, got %d", BurnMessageLen, len(msg.BurnToken))
	}

	msg.Version = binary.BigEndian.Uint32(bz[BurnMsgVersionIndex:BurnTokenIndex])
	msg.BurnToken = bz[BurnTokenIndex:MintRecipientIndex]
	msg.MintRecipient = bz[MintRecipientIndex:AmountIndex]
	msg.Amount = math.NewIntFromBigInt(new(big.Int).SetBytes(bz[AmountIndex:MsgSenderIndex]))
	msg.MessageSender = bz[MsgSenderIndex:BurnMessageLen]

	return msg, nil
}

// Bytes parses a BurnMessage struct into a byte array
// burn token, mint recipient, and message sender must be 32 bytes
func (msg *BurnMessage) Bytes() ([]byte, error) {
	if len(msg.BurnToken) != BurnTokenLen {
		return nil, errors.Wrapf(ErrParsingBurnMessage, "burn token must be 32 bytes, got %d", len(msg.BurnToken))
	}
	if len(msg.MintRecipient) != MintRecipientLen {
		return nil, errors.Wrapf(ErrParsingBurnMessage, "mint recipient must be 32 bytes, got %d", len(msg.MintRecipient))
	}
	if len(msg.MessageSender) != AddressBytesLen {
		return nil, errors.Wrapf(ErrParsingBurnMessage, "message sender must be 32 bytes, got %d", len(msg.MessageSender))
	}

	result := make([]byte, BurnMessageLen)

	versionBytes := make([]byte, VersionLen)
	binary.BigEndian.PutUint32(versionBytes, msg.Version)

	amountBytes := make([]byte, AmountLen)
	msg.Amount.BigInt().FillBytes(amountBytes)

	copy(result[BurnMsgVersionIndex:BurnTokenIndex], versionBytes)
	copy(result[BurnTokenIndex:MintRecipientIndex], msg.BurnToken)
	copy(result[MintRecipientIndex:AmountIndex], msg.MintRecipient)
	copy(result[AmountIndex:MsgSenderIndex], amountBytes)
	copy(result[MsgSenderIndex:BurnMessageLen], msg.MessageSender)

	return result, nil
}
