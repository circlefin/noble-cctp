package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func createNPerMessageBurnLimits(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PerMessageBurnLimit {
	items := make([]types.PerMessageBurnLimit, n)
	for i := range items {
		items[i].Denom = "amount" + strconv.Itoa(i)
		items[i].Amount = math.NewInt(int64(i))
		keeper.SetPerMessageBurnLimit(ctx, items[i])
	}
	return items
}

func TestPerMessageBurnLimit(t *testing.T) {

	keeper, ctx := keepertest.CctpKeeper(t)

	limit, found := keeper.GetPerMessageBurnLimit(ctx, "usdc")
	require.False(t, found)

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "usdc",
		Amount: math.NewInt(123),
	}
	keeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	limit, found = keeper.GetPerMessageBurnLimit(ctx, perMessageBurnLimit.Denom)
	require.True(t, found)
	require.Equal(t,
		perMessageBurnLimit,
		nullify.Fill(&limit),
	)

	newPerMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "usdc",
		Amount: math.NewInt(456),
	}

	keeper.SetPerMessageBurnLimit(ctx, newPerMessageBurnLimit)

	limit, found = keeper.GetPerMessageBurnLimit(ctx, newPerMessageBurnLimit.Denom)
	require.True(t, found)
	require.Equal(t,
		newPerMessageBurnLimit,
		nullify.Fill(&limit),
	)

	separateCurrencyPerMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "euroc",
		Amount: math.NewInt(789),
	}
	keeper.SetPerMessageBurnLimit(ctx, separateCurrencyPerMessageBurnLimit)
	limit, found = keeper.GetPerMessageBurnLimit(ctx, separateCurrencyPerMessageBurnLimit.Denom)

	require.True(t, found)
	require.Equal(t,
		separateCurrencyPerMessageBurnLimit,
		nullify.Fill(&limit),
	)
}

func TestPerMessageBurnLimitsGetAll(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNPerMessageBurnLimits(cctpKeeper, ctx, 10)
	denom := make([]types.PerMessageBurnLimit, len(items))
	for i, item := range items {
		denom[i] = item
	}
	require.ElementsMatch(t,
		nullify.Fill(denom),
		nullify.Fill(cctpKeeper.GetAllPerMessageBurnLimits(ctx)),
	)
}
