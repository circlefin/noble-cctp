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
func TestUnpauseBurningAndMintingHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	message := types.MsgUnpauseBurningAndMinting{
		From: pauser,
	}
	_, err := server.UnpauseBurningAndMinting(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetBurningAndMintingPaused(ctx)
	require.True(t, found)
	require.Equal(t, false, actual.Paused)
}

func TestUnpauseBurningAndMintingAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUnpauseBurningAndMinting{
		From: "authority",
	}
	_, err := server.UnpauseBurningAndMinting(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestUnpauseBurningAndMintingInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	message := types.MsgUnpauseBurningAndMinting{
		From: "not the authority",
	}
	_, err := server.UnpauseBurningAndMinting(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot unpause burning and minting")
}
