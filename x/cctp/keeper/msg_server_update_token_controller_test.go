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

func TestUpdateTokenControllerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	newTokenController := sample.AccAddress()

	message := types.MsgUpdateTokenController{
		From:               owner,
		NewTokenController: newTokenController,
	}

	_, err := server.UpdateTokenController(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual := testkeeper.GetTokenController(ctx)
	require.Equal(t, newTokenController, actual)
}

func TestUpdateTokenControllerAuthorityIsNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgUpdateTokenController{
		From:               "not the authority",
		NewTokenController: sample.AccAddress(),
	}
	require.Panicsf(t, func() {
		_, _ = server.UpdateTokenController(sdk.WrapSDKContext(ctx), &message)
	}, "cctp owner not found in state")
}

func TestUpdateTokenControllerInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	newTokenController := sample.AccAddress()

	message := types.MsgUpdateTokenController{
		From:               sample.AccAddress(),
		NewTokenController: newTokenController,
	}

	_, err := server.UpdateTokenController(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the authority")
}
