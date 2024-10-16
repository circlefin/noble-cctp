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

package keeper

import (
	"bytes"
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) SendMessageWithCaller(goCtx context.Context, msg *types.MsgSendMessageWithCaller) (*types.MsgSendMessageWithCallerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	messageSender := make([]byte, 32)
	fromAccAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	copy(messageSender[12:], fromAccAddress)

	emptyByteArr := make([]byte, types.DestinationCallerLen)
	if len(msg.DestinationCaller) != types.DestinationCallerLen || bytes.Equal(msg.DestinationCaller, emptyByteArr) {
		return nil, errors.Wrap(types.ErrSendMessage, "destination caller must be nonzero")
	}

	nonce := k.ReserveAndIncrementNonce(ctx)

	err = k.sendMessage(
		ctx,
		msg.DestinationDomain,
		msg.Recipient,
		msg.DestinationCaller,
		messageSender,
		nonce.Nonce,
		msg.MessageBody)

	return &types.MsgSendMessageWithCallerResponse{Nonce: nonce.Nonce}, err
}
