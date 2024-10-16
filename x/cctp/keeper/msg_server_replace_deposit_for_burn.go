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
	"encoding/hex"

	"cosmossdk.io/errors"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (k msgServer) ReplaceDepositForBurn(goCtx context.Context, msg *types.MsgReplaceDepositForBurn) (*types.MsgReplaceDepositForBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	paused, found := k.GetBurningAndMintingPaused(ctx)
	if found && paused.Paused {
		return nil, errors.Wrap(types.ErrDepositForBurn, "burning and minting are paused")
	}

	// verify and parse original originalMessage
	originalMessage, err := new(types.Message).Parse(msg.OriginalMessage)
	if err != nil {
		return nil, err
	}

	// verify and parse BurnMessage
	burnMessage, err := new(types.BurnMessage).Parse(originalMessage.MessageBody)
	if err != nil {
		return nil, err
	}

	// validate originalMessage sender is the same as this message sender
	messageSender := make([]byte, 32)
	fromAccAddress, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	copy(messageSender[12:], fromAccAddress)
	if !bytes.Equal(messageSender, burnMessage.MessageSender) {
		return nil, errors.Wrap(types.ErrDepositForBurn, "invalid sender for message")
	}

	// validate new mint recipient
	emptyByteArr := make([]byte, types.MintRecipientLen)
	if bytes.Equal(emptyByteArr, msg.NewMintRecipient) {
		return nil, errors.Wrap(types.ErrDepositForBurn, "mint recipient must be nonzero")
	}

	newMessageBody := types.BurnMessage{
		Version:       burnMessage.Version,
		BurnToken:     burnMessage.BurnToken,
		MintRecipient: msg.NewMintRecipient,
		Amount:        burnMessage.Amount,
		MessageSender: burnMessage.MessageSender,
	}

	newMessageBodyBytes, err := newMessageBody.Bytes()
	if err != nil {
		return nil, errors.Wrapf(types.ErrParsingBurnMessage, "error parsing burn message")
	}

	replaceMessage := types.MsgReplaceMessage{
		From:                 types.ModuleAddress.String(),
		OriginalMessage:      msg.OriginalMessage,
		OriginalAttestation:  msg.OriginalAttestation,
		NewMessageBody:       newMessageBodyBytes,
		NewDestinationCaller: msg.NewDestinationCaller,
	}
	_, err = k.ReplaceMessage(goCtx, &replaceMessage)
	if err != nil {
		return nil, errors.Wrap(err, "error calling replace message")
	}

	event := types.DepositForBurn{
		Nonce:                     originalMessage.Nonce,
		BurnToken:                 hex.EncodeToString(crypto.Keccak256(burnMessage.BurnToken)),
		Amount:                    burnMessage.Amount,
		Depositor:                 msg.From,
		MintRecipient:             msg.NewMintRecipient,
		DestinationDomain:         originalMessage.DestinationDomain,
		DestinationTokenMessenger: originalMessage.Recipient,
		DestinationCaller:         msg.NewDestinationCaller,
	}

	err = ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgReplaceDepositForBurnResponse{}, err
}
