package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func TestAuthority(t *testing.T) {

	keeper, ctx := keepertest.CctpKeeper(t)

	authority, found := keeper.GetAuthority(ctx)
	require.False(t, found)

	Authority := types.Authority{Address: "1"}
	keeper.SetAuthority(ctx, Authority)

	authority, found = keeper.GetAuthority(ctx)
	require.True(t, found)
	require.Equal(t,
		Authority,
		nullify.Fill(&authority),
	)

	newAuthority := types.Authority{Address: "2"}

	keeper.SetAuthority(ctx, newAuthority)

	authority, found = keeper.GetAuthority(ctx)
	require.True(t, found)
	require.Equal(t,
		newAuthority,
		nullify.Fill(&authority),
	)
}
