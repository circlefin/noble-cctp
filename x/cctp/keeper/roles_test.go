package keeper_test

import (
	"testing"

	"github.com/circlefin/noble-cctp/testutil/sample"

	"github.com/stretchr/testify/require"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
)

func TestOwner(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)

	owner := sample.AccAddress()
	keeper.SetOwner(ctx, owner)

	require.Equal(t, owner, keeper.GetOwner(ctx))

	newOwner := sample.AccAddress()
	keeper.SetOwner(ctx, newOwner)

	require.Equal(t, newOwner, keeper.GetOwner(ctx))
}
