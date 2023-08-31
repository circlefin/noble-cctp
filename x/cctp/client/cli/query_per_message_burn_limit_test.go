/*
 * Copyright (c) 2023, © Circle Internet Financial, LTD.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"google.golang.org/grpc/codes"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/circlefin/noble-cctp/testutil/network"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/client/cli"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func networkWithPerMessageBurnLimitObjects(t *testing.T, n int) (*network.Network, []types.PerMessageBurnLimit) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		limit := types.PerMessageBurnLimit{
			Denom:  strconv.Itoa(i),
			Amount: math.NewInt(int64(i)),
		}
		nullify.Fill(&limit)
		state.PerMessageBurnLimitList = append(state.PerMessageBurnLimitList, limit)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.PerMessageBurnLimitList
}

func TestShowPerMessageBurnLimit(t *testing.T) {
	net, objs := networkWithPerMessageBurnLimitObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc  string
		denom string

		args []string
		err  error
		obj  types.PerMessageBurnLimit
	}{
		{
			desc:  "found",
			denom: objs[0].Denom,
			args:  common,
			obj:   objs[0],
		},
		{
			desc:  "not found",
			denom: "123",
			args:  common,
			err:   status.Error(codes.NotFound, "not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.denom,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowPerMessageBurnLimit(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetPerMessageBurnLimitResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.BurnLimit)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.BurnLimit),
				)
			}
		})
	}
}

func TestListPerMessageBurnLimits(t *testing.T) {
	net, objs := networkWithPerMessageBurnLimitObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerMessageBurnLimits(), args)
			require.NoError(t, err)
			var resp types.QueryAllPerMessageBurnLimitsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.BurnLimits), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.BurnLimits),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerMessageBurnLimits(), args)
			require.NoError(t, err)
			var resp types.QueryAllPerMessageBurnLimitsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.BurnLimits), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.BurnLimits),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerMessageBurnLimits(), args)
		require.NoError(t, err)
		var resp types.QueryAllPerMessageBurnLimitsResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.BurnLimits),
		)
	})
}
