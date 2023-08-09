package keeper_test

import (
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
* Happy path
* Authority not set
* Invalid authority
* Token messenger not found
 */

func TestRemoveTokenMessengerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	addMessage := types.MsgAddTokenMessenger{
		From:     authority.Address,
		DomainId: 16,
		Address:  "address to remove",
	}

	_, err := server.AddTokenMessenger(sdk.WrapSDKContext(ctx), &addMessage)
	require.Nil(t, err)

	removeMessage := types.MsgRemoveTokenMessenger{
		From:     authority.Address,
		DomainId: addMessage.DomainId,
	}

	_, err = server.RemoveTokenMessenger(sdk.WrapSDKContext(ctx), &removeMessage)
	require.Nil(t, err)

	_, found := testkeeper.GetTokenMessenger(ctx, removeMessage.DomainId)
	require.False(t, found)

}

func TestRemoveTokenMessengerAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgRemoveTokenMessenger{
		From:     sample.AccAddress(),
		DomainId: 16,
	}

	_, err := server.RemoveTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestRemoveTokenMessengerInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgRemoveTokenMessenger{
		From:     "not the authority address",
		DomainId: 16,
	}

	_, err := server.RemoveTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot remove token messengers")
}

func TestRemoveTokenMessengerTokenMessengerNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgRemoveTokenMessenger{
		From:     authority.Address,
		DomainId: 1,
	}

	_, err := server.RemoveTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrTokenMessengerNotFound, err)
	require.Contains(t, err.Error(), "no token messenger was found for this domain")
}
