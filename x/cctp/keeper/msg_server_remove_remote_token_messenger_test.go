package keeper_test

import (
	"testing"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

/*
* Happy path
* Authority not set
* Invalid authority
* Remote token messenger not found
 */

func TestRemoveRemoteTokenMessengerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	addMessage := types.MsgAddRemoteTokenMessenger{
		From:     owner,
		DomainId: 16,
		Address:  "address to remove",
	}

	_, err := server.AddRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &addMessage)
	require.Nil(t, err)

	removeMessage := types.MsgRemoveRemoteTokenMessenger{
		From:     owner,
		DomainId: addMessage.DomainId,
	}

	_, err = server.RemoveRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &removeMessage)
	require.Nil(t, err)

	_, found := testkeeper.GetRemoteTokenMessenger(ctx, removeMessage.DomainId)
	require.False(t, found)
}

func TestRemoveRemoteTokenMessengerAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgRemoveRemoteTokenMessenger{
		From:     sample.AccAddress(),
		DomainId: 16,
	}

	_, err := server.RemoveRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestRemoveRemoteTokenMessengerInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgRemoveRemoteTokenMessenger{
		From:     "not the authority address",
		DomainId: 16,
	}

	_, err := server.RemoveRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot remove remote token messengers")
}

func TestRemoveRemoteTokenMessengerTokenMessengerNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgRemoveRemoteTokenMessenger{
		From:     owner,
		DomainId: 1,
	}

	_, err := server.RemoveRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrRemoteTokenMessengerNotFound, err)
	require.Contains(t, err.Error(), "no remote token messenger was found for this domain")
}
