package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/require"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func TestPerMessageBurnLimitQuery(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)

	perMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(int64(21)),
	}
	keeper.SetPerMessageBurnLimit(ctx, perMessageBurnLimit)

	rst, found := keeper.GetPerMessageBurnLimit(ctx, perMessageBurnLimit.Denom)
	require.True(t, found)
	require.Equal(t,
		perMessageBurnLimit,
		nullify.Fill(&rst),
	)

	newPerMessageBurnLimit := types.PerMessageBurnLimit{
		Denom:  "uusdc",
		Amount: math.NewInt(int64(22)),
	}

	keeper.SetPerMessageBurnLimit(ctx, newPerMessageBurnLimit)

	rst, found = keeper.GetPerMessageBurnLimit(ctx, newPerMessageBurnLimit.Denom)
	require.True(t, found)
	require.Equal(t,
		newPerMessageBurnLimit,
		nullify.Fill(&rst),
	)
}

func TestPerMessageBurnLimitQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPerMessageBurnLimits(keeper, ctx, 5)
	perMessageBurnLimits := make([]types.PerMessageBurnLimit, len(msgs))
	copy(perMessageBurnLimits, msgs)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPerMessageBurnLimitsRequest {
		return &types.QueryAllPerMessageBurnLimitsRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(perMessageBurnLimits); i += step {
			resp, err := keeper.PerMessageBurnLimits(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.BurnLimits), step)
			require.Subset(t,
				nullify.Fill(perMessageBurnLimits),
				nullify.Fill(resp.BurnLimits),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(perMessageBurnLimits); i += step {
			resp, err := keeper.PerMessageBurnLimits(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.BurnLimits), step)
			require.Subset(t,
				nullify.Fill(perMessageBurnLimits),
				nullify.Fill(resp.BurnLimits),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.PerMessageBurnLimits(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(perMessageBurnLimits), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(perMessageBurnLimits),
			nullify.Fill(resp.BurnLimits),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.Attesters(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
