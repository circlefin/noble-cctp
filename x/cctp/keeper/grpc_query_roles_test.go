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

func TestRolesQuery(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)

	owner := "test-owner"
	keeper.SetOwner(ctx, owner)
	attesterManager := "test-attester-manager"
	keeper.SetAttesterManager(ctx, attesterManager)
	pauser := "test-pauser"
	keeper.SetPauser(ctx, pauser)
	tokenController := "test-token-controller"
	keeper.SetTokenController(ctx, tokenController)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryRolesRequest
		response *types.QueryRolesResponse
		err      error
	}{
		{
			desc:    "First",
			request: &types.QueryRolesRequest{},
			response: &types.QueryRolesResponse{
				Owner:           owner,
				AttesterManager: attesterManager,
				Pauser:          pauser,
				TokenController: tokenController,
			},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Roles(goCtx, tc.request)
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
