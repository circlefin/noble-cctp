package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func TestMaxMessageBodySize(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)

	_, found := keeper.GetMaxMessageBodySize(ctx)
	require.False(t, found)

	MaxMessageBodySize := types.MaxMessageBodySize{Amount: 21}
	keeper.SetMaxMessageBodySize(ctx, MaxMessageBodySize)

	maxMessageBodySize, found := keeper.GetMaxMessageBodySize(ctx)
	require.True(t, found)
	require.Equal(t,
		MaxMessageBodySize,
		nullify.Fill(&maxMessageBodySize),
	)

	newMaxMessageBodySize := types.MaxMessageBodySize{Amount: 22}

	keeper.SetMaxMessageBodySize(ctx, newMaxMessageBodySize)

	maxMessageBodySize, found = keeper.GetMaxMessageBodySize(ctx)
	require.True(t, found)
	require.Equal(t,
		newMaxMessageBodySize,
		nullify.Fill(&maxMessageBodySize),
	)
}
