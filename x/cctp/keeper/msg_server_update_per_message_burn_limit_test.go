package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
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

func TestUpdatePerMessageBurnLimitHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgUpdatePerMessageBurnLimit{
		From:   tokenController,
		Denom:  "uusdc",
		Amount: math.NewInt(123),
	}

	_, err := server.UpdatePerMessageBurnLimit(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetPerMessageBurnLimit(ctx, message.Denom)
	require.True(t, found)
	require.Equal(t, message.Denom, actual.Denom)
	require.Equal(t, message.Amount, actual.Amount)
}

func TestUpdatePerMessageBurnLimitAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUpdatePerMessageBurnLimit{
		From:   sample.AccAddress(),
		Denom:  "uusdc",
		Amount: math.NewInt(123),
	}

	_, err := server.UpdatePerMessageBurnLimit(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestUpdatePerMessageBurnLimitInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgUpdatePerMessageBurnLimit{
		From:   "not authority",
		Denom:  "uusdc",
		Amount: math.NewInt(123),
	}

	_, err := server.UpdatePerMessageBurnLimit(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the per message burn limit")
}
