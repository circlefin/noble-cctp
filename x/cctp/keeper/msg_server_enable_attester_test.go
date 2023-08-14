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
 * Attester already found
 */
func TestEnableAttesterHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgEnableAttester{
		From:     attesterManager,
		Attester: []byte("attester"),
	}

	_, err := server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetAttester(ctx, string(message.Attester))
	require.True(t, found)
	require.Equal(t, message.Attester, []byte(actual.Attester))
}

func TestEnableAttesterAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgEnableAttester{
		From:     sample.AccAddress(),
		Attester: []byte("attester"),
	}

	_, err := server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestEnableAttesterInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgEnableAttester{
		From:     sample.AccAddress(),
		Attester: []byte("attester"),
	}

	_, err := server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot enable attesters")
}

func TestEnableAttesterAttesterAlreadyFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existingAttester := types.Attester{Attester: "attester"}
	testkeeper.SetAttester(ctx, existingAttester)

	message := types.MsgEnableAttester{
		From:     attesterManager,
		Attester: []byte("attester"),
	}

	_, err := server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAttesterAlreadyFound, err)
	require.Contains(t, err.Error(), "this attester already exists in the store")
}
