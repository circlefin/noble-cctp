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
package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	fiattokenfactorytypes "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

/*
* Happy path
* 0 amount
* Negative amount
* Nil mint recipient
* Empty mint recipient
* Remote Token Messenger not found
* Minting Denom not found
* Burning and Minting is paused
* Amount is greater than per message burn limit
* Burn fails
 */

func TestDepositForBurnHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	fiatfkeeper, fiatfctx := keepertest.FiatTokenfactoryKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	fiatfkeeper.SetMintingDenom(fiatfctx, fiattokenfactorytypes.MintingDenom{Denom: "uUsDC"})

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurn{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
	}
	resp, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.Nil(t, err)
	require.Equal(t, startingNonce.Nonce, resp.Nonce)

	nextNonce, found := testkeeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t, startingNonce.Nonce+1, nextNonce.Nonce)
}

func TestDepositForBurnZeroAmount(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(0),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
	}
	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "amount must be positive")
}

func TestDepositForBurnNegativeAmount(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(-59),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
	}
	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "amount must be positive")
}

func TestDepositForBurnEmptyMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     make([]byte, types.MintRecipientLen),
		BurnToken:         "uUsDC",
	}
	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}

func TestDepositForBurnNilMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		BurnToken:         "uUsDC",
	}
	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}

func TestDepositForBurnTokenMessengerNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	msg := types.MsgDepositForBurn{
		From:              "sender address",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
	}
	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "unable to look up destination token messenger")
}

func TestDepositForBurnMintingDenomNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "not usdc",
	}

	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "is not supported")
}

func TestDepositForBurnBurningAndMintingIsPaused(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	testkeeper.SetBurningAndMintingPaused(ctx, types.BurningAndMintingPaused{Paused: true})

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
	}

	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "burning and minting are paused")
}

func TestDepositForBurnAmountIsGreaterThanPerMessageBurnLimit(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
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

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
	}

	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "cannot burn more than the maximum per message burn limit")
}

func TestDepositForBurnBurnFails(t *testing.T) {
	testkeeper, ctx := keepertest.ErrCctpKeeper(t)
	fiattfkeeper, fiattfctx := keepertest.ErrFiatTokenfactoryKeeper(t)

	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)
	fiattfkeeper.SetMintingDenom(fiattfctx, fiattokenfactorytypes.MintingDenom{Denom: "uUsDC"})

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurn{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
	}
	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.Contains(t, err.Error(), "tokens can not be burned")
}

func TestDepositForBurnTransferFails(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeperWithErrBank(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(1_000_000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurn{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(42),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
	}

	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorContains(t, err, "error during transfer: intentional error")
}

func TestDepositForBurnMessageFormatFails(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	fiatfkeeper, fiatfctx := keepertest.FiatTokenfactoryKeeper(t)

	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)
	fiatfkeeper.SetMintingDenom(fiatfctx, fiattokenfactorytypes.MintingDenom{Denom: "uUsDC"})

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurn{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     common.FromHex("0xfCE4cE85e1F74C01e0ecccd8BbC4606f83D3FC90"),
		BurnToken:         "uUsDC",
	}

	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, err, types.ErrParsingBurnMessage)
}

func TestDepositForBurnSendMessageFails(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	fiatfkeeper, fiatfctx := keepertest.FiatTokenfactoryKeeper(t)

	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	fiatfkeeper.SetMintingDenom(fiatfctx, fiattokenfactorytypes.MintingDenom{Denom: "uUsDC"})

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	testkeeper.SetSendingAndReceivingMessagesPaused(ctx, types.SendingAndReceivingMessagesPaused{Paused: true})

	msg := types.MsgDepositForBurn{
		From:              sample.AccAddress(),
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
	}

	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrSendMessage, err)
}
