package cli_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"google.golang.org/grpc/codes"
	"strconv"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/circlefin/noble-cctp/x/cctp/client/cli"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/strangelove-ventures/noble/testutil/network"
	"github.com/strangelove-ventures/noble/testutil/nullify"
)

func networkWithTokenMessengerObjects(t *testing.T, n int) (*network.Network, []types.TokenMessenger) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		tokenMessenger := types.TokenMessenger{
			DomainId: uint32(i),
			Address:  strconv.Itoa(i),
		}
		nullify.Fill(&tokenMessenger)
		state.TokenMessengerList = append(state.TokenMessengerList, tokenMessenger)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.TokenMessengerList
}

func TestShowTokenMessenger(t *testing.T) {
	net, objs := networkWithTokenMessengerObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc         string
		remoteDomain string
		remoteToken  string

		args []string
		err  error
		obj  types.TokenMessenger
	}{
		{
			desc:         "found",
			remoteDomain: strconv.Itoa(int(objs[0].DomainId)),
			remoteToken:  objs[0].Address,
			args:         common,
			obj:          objs[0],
		},
		{
			desc:         "not found",
			remoteDomain: "notakey",
			remoteToken:  objs[0].Address,
			args:         common,
			err:          status.Error(codes.NotFound, "not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.remoteDomain,
				tc.remoteToken,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowTokenMessenger(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetTokenMessengerResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.TokenMessenger.DomainId)
				require.NotNil(t, resp.TokenMessenger.Address)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.TokenMessenger),
				)
			}
		})
	}
}

func TestListTokenMessengers(t *testing.T) {
	net, objs := networkWithTokenMessengerObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTokenMessengers(), args)
			require.NoError(t, err)
			var resp types.QueryAllTokenMessengersResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.TokenMessengers), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.TokenMessengers),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTokenMessengers(), args)
			require.NoError(t, err)
			var resp types.QueryAllTokenMessengersResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.TokenMessengers), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.TokenMessengers),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTokenMessengers(), args)
		require.NoError(t, err)
		var resp types.QueryAllTokenMessengersResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.TokenMessengers),
		)
	})
}
