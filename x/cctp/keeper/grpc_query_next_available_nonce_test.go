package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func TestNextAvailableNonceQuery(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	nonce := types.Nonce{Nonce: uint64(123)}
	keeper.SetNextAvailableNonce(ctx, nonce)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetNextAvailableNonceRequest
		response *types.QueryGetNextAvailableNonceResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetNextAvailableNonceRequest{},
			response: &types.QueryGetNextAvailableNonceResponse{Nonce: nonce},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.NextAvailableNonce(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
