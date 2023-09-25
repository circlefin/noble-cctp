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

func TestUpdateAttesterManagerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	newAttesterManager := sample.AccAddress()

	message := types.MsgUpdateAttesterManager{
		From:               owner,
		NewAttesterManager: newAttesterManager,
	}

	_, err := server.UpdateAttesterManager(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual := testkeeper.GetAttesterManager(ctx)
	require.Equal(t, newAttesterManager, actual)
}

func TestUpdateAttesterManagerAuthorityIsNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUpdateAttesterManager{
		From:               "not the authority",
		NewAttesterManager: sample.AccAddress(),
	}
	require.Panicsf(t, func() {
		_, _ = server.UpdateAttesterManager(sdk.WrapSDKContext(ctx), &message)
	}, "cctp owner not found in state")
}

func TestUpdateAttesterManagerInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	newAttesterManager := sample.AccAddress()

	message := types.MsgUpdateAttesterManager{
		From:               sample.AccAddress(),
		NewAttesterManager: newAttesterManager,
	}

	_, err := server.UpdateAttesterManager(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the attester manager")
}
