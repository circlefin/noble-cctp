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

package keeper_test

import (
	"testing"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

/*
 * Happy path
 * Sending and receiving messages is paused
 * Message body is too long
 * Recipient is empty
 * Invalid from address
 */
func TestSendMessageHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessage{
		From:              sample.AccAddress(),
		DestinationDomain: 3,
		Recipient:         []byte("12345678901234567890123456789012"),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
	}

	resp, err := server.SendMessage(ctx, &msg)
	require.Nil(t, err)
	require.Equal(t, nonce.Nonce, resp.Nonce)

	nextNonce, found := testkeeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t, nonce.Nonce+1, nextNonce.Nonce)
}

func TestSendMessageSendingAndReceivingMessagesPaused(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: true}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	_, err := server.SendMessage(ctx, &types.MsgSendMessage{
		From: sample.AccAddress(),
	})
	require.ErrorIs(t, types.ErrSendMessage, err)
	require.Contains(t, err.Error(), "sending and receiving messages is paused")
}

func TestSendMessageMessageBodyTooLong(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	max := types.MaxMessageBodySize{Amount: 5}
	testkeeper.SetMaxMessageBodySize(ctx, max)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessage{
		From:              sample.AccAddress(),
		DestinationDomain: 3,
		Recipient:         []byte("12345678901234567890123456789012"),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
	}

	_, err := server.SendMessage(ctx, &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
	require.Contains(t, err.Error(), "message body exceeds max size")
}

func TestSendMessageRecipientEmpty(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessage{
		From:              sample.AccAddress(),
		DestinationDomain: 3,
		Recipient:         make([]byte, types.MintRecipientLen),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
	}

	_, err := server.SendMessage(ctx, &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
	require.Contains(t, err.Error(), "recipient must not be nonzero")
}

func TestSendMessageRecipientInvalid(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgSendMessage{
		From:      sample.AccAddress(),
		Recipient: common.FromHex("0xfCE4cE85e1F74C01e0ecccd8BbC4606f83D3FC90"),
	}

	_, err := server.SendMessage(ctx, &msg)
	require.ErrorIs(t, err, types.ErrParsingMessage)
}

func TestSendMessageInvalidSenderAddress(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	paused := types.SendingAndReceivingMessagesPaused{Paused: false}
	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, paused)

	nonce := types.Nonce{
		Nonce: 5,
	}
	testkeeper.SetNextAvailableNonce(ctx, nonce)

	msg := types.MsgSendMessage{
		From:              "invalid sender address",
		DestinationDomain: 3,
		Recipient:         []byte("12345678901234567890123456789012"),
		MessageBody:       []byte("It's not about money, it's about sending a message"),
	}

	_, err := server.SendMessage(ctx, &msg)
	require.ErrorIs(t, err, types.ErrInvalidAddress)
	require.Contains(t, err.Error(), "invalid from address")
}
