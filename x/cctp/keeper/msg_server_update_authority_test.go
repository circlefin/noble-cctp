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

func TestUpdateAuthorityHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: "current address"}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgUpdateAuthority{
		From:         authority.Address,
		NewAuthority: "new address",
	}
	_, err := server.UpdateAuthority(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetAuthority(ctx)
	require.True(t, found)
	require.Equal(t, message.NewAuthority, actual.Address)
}

func TestUpdateAuthorityAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUpdateAuthority{
		From:         sample.AccAddress(),
		NewAuthority: "new address",
	}
	_, err := server.UpdateAuthority(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestUpdateAuthorityInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: "address"}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgUpdateAuthority{
		From:         "not the authority",
		NewAuthority: "new address",
	}
	_, err := server.UpdateAuthority(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the authority")
}
