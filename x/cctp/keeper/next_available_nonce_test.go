package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func TestNextAvailableNonce(t *testing.T) {

	keeper, ctx := keepertest.CctpKeeper(t)

	next, found := keeper.GetNextAvailableNonce(ctx)
	require.False(t, found)

	savedNonce := types.Nonce{Nonce: 21}
	keeper.SetNextAvailableNonce(ctx, savedNonce)

	next, found = keeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t,
		savedNonce,
		nullify.Fill(&next),
	)

	newSavedNonce := types.Nonce{Nonce: 22}

	keeper.SetNextAvailableNonce(ctx, newSavedNonce)

	next, found = keeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t,
		newSavedNonce,
		nullify.Fill(&next),
	)
}
