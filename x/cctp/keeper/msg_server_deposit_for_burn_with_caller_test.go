package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	fiattokenfactorytypes "github.com/strangelove-ventures/noble/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
)

/*
* Happy path
* Invalid destination caller
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
	testkeeper, ctx := keepertest.CctpKeeper(t)
	fiatfkeeper, fiatfctx := keepertest.FiatTokenfactoryKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  "12345678901234567890123456789012",
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	fiatfkeeper.SetMintingDenom(fiatfctx, fiattokenfactorytypes.MintingDenom{Denom: "uUsDC"})

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	resp, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.Nil(t, err)
	require.Equal(t, startingNonce.Nonce, resp.Nonce)

	nextNonce, found := testkeeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t, startingNonce.Nonce+1, nextNonce.Nonce)
}

func TestDepositForBurnWithCallerInvalidDestinationCaller(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(0),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: make([]byte, types.DestinationCallerLen),
	}
	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrInvalidDestinationCaller, err)
	require.Contains(t, err.Error(), "invalid destination caller")
}

func TestDepositForBurnWithCallerZeroAmount(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(0),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "amount must be positive")
}

func TestDepositForBurnWithCallerNegativeAmount(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(-4738),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "amount must be positive")
}

func TestDepositForBurnWithCallerNilMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     make([]byte, types.MintRecipientLen),
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}

func TestDepositForBurnWithCallerEmptyMintRecipient(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		BurnToken:         "uUsDC",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "mint recipient must be nonzero")
}

func TestDepositForBurnWithCallerTokenMessengerNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender address",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}
	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrDepositForBurn, err)
	require.Contains(t, err.Error(), "unable to look up destination token messenger")
}

func TestDepositForBurnWithCallerMintingDenomNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  "destination-remote-token-messenger",
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "not usdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "is not supported")
}

func TestDepositForBurnWithCallerBurningAndMintingIsPaused(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  "destination-remote-token-messenger",
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	testkeeper.SetBurningAndMintingPaused(ctx, types.BurningAndMintingPaused{Paused: true})

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "burning and minting are paused")
}

func TestDepositForBurnWithCallerAmountIsGreaterThanPerMessageBurnLimit(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	startingNonce := types.Nonce{Nonce: 1}
	testkeeper.SetNextAvailableNonce(ctx, startingNonce)

	remoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  "destination-remote-token-messenger",
	}
	testkeeper.SetRemoteTokenMessenger(ctx, remoteTokenMessenger)

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(4),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurnWithCaller{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uusdc",
		DestinationCaller: []byte("12345678901234567890123456789012"),
	}

	_, err := server.DepositForBurnWithCaller(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "cannot burn more than the maximum per message burn limit")
}
