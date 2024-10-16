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

	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/stretchr/testify/require"
)

/*
* Happy path
* Invalid destination caller
* Invalid from address
* 0 amount
* Negative amount
* Nil mint recipient
* Empty mint recipient
* Remote Token Messenger not found
* Minting Denom not found
* Burning and Minting is paused
* Amount is greater than per message burn limit
 */

func TestDepositForBurnWithCallerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	resp, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.Nil(t, err)
	require.Equal(t, startingNonce.Nonce, resp.Nonce)

	nextNonce, found := testkeeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t, startingNonce.Nonce+1, nextNonce.Nonce)
}

func TestDepositForBurnWithCallerInvalidDestinationCaller(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(0),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: make([]byte, types.DestinationCallerLen),
	}
	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrInvalidDestinationCaller, err)
	require.Contains(t, err.Error(), "invalid destination caller")
}
func TestDepositForBurnWithCallerInvalidFromAddress(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "invalid from address",
		Amount:            math.NewInt(0),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, err, types.ErrInvalidAddress)
	require.Contains(t, err.Error(), "invalid from address")
}

func TestDepositForBurnWithCallerZeroAmount(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(0),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "amount must be positive")
}

func TestDepositForBurnWithCallerNegativeAmount(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(-4738),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "amount must be positive")
}

func TestDepositForBurnWithCallerNilMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     make([]byte, types.MintRecipientLen),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}

func TestDepositForBurnWithCallerEmptyMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}

func TestDepositForBurnWithCallerTokenMessengerNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "unable to look up destination token messenger")
}

func TestDepositForBurnWithCallerMintingDenomNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "not usdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "is not supported")
}

func TestDepositForBurnWithCallerBurningAndMintingIsPaused(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	testkeeper.SetBurningAndMintingPaused(ctx, types.BurningAndMintingPaused{Paused: true})

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "burning and minting are paused")
}

func TestDepositForBurnWithCallerAmountIsGreaterThanPerMessageBurnLimit(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(4),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "cannot burn more than the maximum per message burn limit")
}

func TestDepositForBurnWithCallerSendMessageFails(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, types.SendingAndReceivingMessagesPaused{Paused: true})

	msg := types.MsgDepositForBurnWithCaller{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.DepositForBurnWithCaller(ctx, &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
}
