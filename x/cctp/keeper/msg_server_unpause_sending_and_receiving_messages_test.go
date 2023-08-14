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
 */
func TestUnpauseSendingAndReceivingMessagesHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	message := types.MsgUnpauseSendingAndReceivingMessages{
		From: pauser,
	}
	_, err := server.UnpauseSendingAndReceivingMessages(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetSendingAndReceivingMessagesPaused(ctx)
	require.True(t, found)
	require.Equal(t, false, actual.Paused)
}

func TestUnpauseSendingAndReceivingMessagesAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUnpauseSendingAndReceivingMessages{
		From: "authority",
	}
	_, err := server.UnpauseSendingAndReceivingMessages(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestUnpauseSendingAndReceivingMessagesInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	message := types.MsgUnpauseSendingAndReceivingMessages{
		From: "not the authority",
	}
	_, err := server.UnpauseSendingAndReceivingMessages(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot unpause sending and receiving messages")
}
