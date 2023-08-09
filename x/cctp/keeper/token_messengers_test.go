package keeper_test

import (
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"strconv"
	"testing"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNTokenMessengers(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.TokenMessenger {
	items := make([]types.TokenMessenger, n)
	for i := range items {
		items[i].DomainId = uint32(i)
		items[i].Address = strconv.Itoa(i)

		keeper.SetTokenMessenger(ctx, items[i])
	}
	return items
}

func TestTokenMessengersGet(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNTokenMessengers(cctpKeeper, ctx, 10)
	for _, item := range items {
		tokenMessenger, found := cctpKeeper.GetTokenMessenger(ctx, item.DomainId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&tokenMessenger),
		)
	}
}

func TestTokenMessengersRemove(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNTokenMessengers(cctpKeeper, ctx, 10)
	for _, item := range items {
		cctpKeeper.DeleteTokenMessenger(ctx, item.DomainId)
		_, found := cctpKeeper.GetTokenMessenger(ctx, item.DomainId)
		require.False(t, found)
	}
}

func TestTokenMessengersGetAll(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNTokenMessengers(cctpKeeper, ctx, 10)
	denom := make([]types.TokenMessenger, len(items))
	for i, item := range items {
		denom[i] = item
	}
	require.ElementsMatch(t,
		nullify.Fill(denom),
		nullify.Fill(cctpKeeper.GetAllTokenMessengers(ctx)),
	)
}
