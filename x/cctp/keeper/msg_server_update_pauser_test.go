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

func TestUpdatePauserHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	newPauser := sample.AccAddress()

	message := types.MsgUpdatePauser{
		From:      owner,
		NewPauser: newPauser,
	}

	_, err := server.UpdatePauser(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual := testkeeper.GetPauser(ctx)
	require.Equal(t, newPauser, actual)
}

func TestUpdatePauserAuthorityIsNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	message := types.MsgUpdatePauser{
		From:      "not the authority",
		NewPauser: sample.AccAddress(),
	}
	require.Panicsf(t, func() {
		_, _ = server.UpdatePauser(sdk.WrapSDKContext(ctx), &message)
	}, "cctp owner not found in state")
}

func TestUpdatePauserInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	newPauser := sample.AccAddress()

	message := types.MsgUpdatePauser{
		From:      sample.AccAddress(),
		NewPauser: newPauser,
	}

	_, err := server.UpdatePauser(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the pauser")
}
