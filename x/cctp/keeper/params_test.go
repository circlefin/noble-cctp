package keeper_test

import (
	"testing"

	testkeeper "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.CctpKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
