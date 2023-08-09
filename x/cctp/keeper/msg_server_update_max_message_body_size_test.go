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
 */

func TestUpdateMaxMessageBodySizeHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: "current address"}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgUpdateMaxMessageBodySize{
		From:        authority.Address,
		MessageSize: uint64(1023),
	}
	_, err := server.UpdateMaxMessageBodySize(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetMaxMessageBodySize(ctx)
	require.True(t, found)
	require.Equal(t, message.MessageSize, actual.Amount)
}

func TestUpdateMaxMessageBodySizeAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUpdateMaxMessageBodySize{
		From:        sample.AccAddress(),
		MessageSize: uint64(1023),
	}
	_, err := server.UpdateMaxMessageBodySize(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestUpdateMaxMessageBodySizeInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: "current address"}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgUpdateMaxMessageBodySize{
		From:        "not the authority",
		MessageSize: uint64(1023),
	}
	_, err := server.UpdateMaxMessageBodySize(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the max message body size")
}
