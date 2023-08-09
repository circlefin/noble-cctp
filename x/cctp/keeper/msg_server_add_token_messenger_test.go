package keeper_test

import (
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/noble/testutil/sample"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
* Happy path
* Authority not set
* Invalid authority
* Token messenger already found
 */

func TestAddTokenMessengerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgAddTokenMessenger{
		From:     authority.Address,
		DomainId: 16,
		Address:  "token_messenger_address",
	}

	_, err := server.AddTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetTokenMessenger(ctx, message.DomainId)
	require.True(t, found)

	require.Equal(t, message.DomainId, actual.DomainId)
	require.Equal(t, message.Address, actual.Address)

}

func TestAddTokenMessengerAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgAddTokenMessenger{
		From:     sample.AccAddress(),
		DomainId: 16,
		Address:  "token_messenger_address",
	}

	_, err := server.AddTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestAddTokenMessengerInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgAddTokenMessenger{
		From:     "not the authority address",
		DomainId: 16,
		Address:  "token_messenger_address",
	}

	_, err := server.AddTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot add token messengers")
}

func TestAddTokenMessengerTokenMessengerAlreadyFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	existingTokenMessenger := types.TokenMessenger{
		DomainId: 3,
		Address:  sample.AccAddress(),
	}
	testkeeper.SetTokenMessenger(ctx, existingTokenMessenger)

	message := types.MsgAddTokenMessenger{
		From:     authority.Address,
		DomainId: existingTokenMessenger.DomainId,
		Address:  "token_messenger_address",
	}

	_, err := server.AddTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrTokenMessengerAlreadyFound, err)
	require.Contains(t, err.Error(), "a token messenger for this domain already exists")
}
