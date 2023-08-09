package keeper_test

import (
	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	fiattokenfactorytypes "github.com/strangelove-ventures/noble/x/fiattokenfactory/types"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
* Happy path
* 0 amount
* Negative amount
* Nil mint recipient
* Empty mint recipient
* Token Messenger not found
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

	tokenMessenger := types.TokenMessenger{
		DomainId: 0,
		Address:  "destination-token-messenger",
	}
	testkeeper.SetTokenMessenger(ctx, tokenMessenger)

	fiatfkeeper.SetMintingDenom(fiatfctx, fiattokenfactorytypes.MintingDenom{Denom: "uUsDC"})

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
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

	tokenMessenger := types.TokenMessenger{
		DomainId: 0,
		Address:  "destination-token-messenger",
	}
	testkeeper.SetTokenMessenger(ctx, tokenMessenger)

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

	tokenMessenger := types.TokenMessenger{
		DomainId: 0,
		Address:  "destination-token-messenger",
	}
	testkeeper.SetTokenMessenger(ctx, tokenMessenger)

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

	tokenMessenger := types.TokenMessenger{
		DomainId: 0,
		Address:  "destination-token-messenger",
	}
	testkeeper.SetTokenMessenger(ctx, tokenMessenger)

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

	tokenMessenger := types.TokenMessenger{
		DomainId: 0,
		Address:  "destination-token-messenger",
	}
	testkeeper.SetTokenMessenger(ctx, tokenMessenger)
	fiattfkeeper.SetMintingDenom(fiattfctx, fiattokenfactorytypes.MintingDenom{Denom: "uUsDC"})

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(800000),
	}
	testkeeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	msg := types.MsgDepositForBurn{
		From:              "sender-address567890123456789012",
		Amount:            math.NewInt(531),
		DestinationDomain: 0,
		MintRecipient:     []byte("12345678901234567890123456789012"),
		BurnToken:         "uUsDC",
	}
	_, err := server.DepositForBurn(sdk.WrapSDKContext(ctx), &msg)
	require.ErrorIs(t, types.ErrBurn, err)
	require.Contains(t, err.Error(), "tokens can not be burned")
}
