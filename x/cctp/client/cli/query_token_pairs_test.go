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
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

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

func networkWithTokenPairObjects(t *testing.T, n int) (*network.Network, []types.TokenPair) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		tokenPair := types.TokenPair{
			RemoteDomain: uint32(i),
			RemoteToken:  []byte{byte(i)},
			LocalToken:   strconv.Itoa(i),
		}
		nullify.Fill(&tokenPair)
		state.TokenPairList = append(state.TokenPairList, tokenPair)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.TokenPairList
}

func TestShowTokenPair(t *testing.T) {
	net, objs := networkWithTokenPairObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc         string
		remoteDomain string
		remoteToken  []byte

		args []string
		err  error
		obj  types.TokenPair
	}{
		{
			desc:         "found",
			remoteDomain: strconv.Itoa(int(objs[0].RemoteDomain)),
			remoteToken:  objs[0].RemoteToken,
			args:         common,
			obj:          objs[0],
		},
		{
			desc:         "not found",
			remoteDomain: "notakey",
			remoteToken:  objs[0].RemoteToken,
			args:         common,
			err:          status.Error(codes.NotFound, "not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.remoteDomain,
				hex.EncodeToString(tc.remoteToken),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowTokenPair(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetTokenPairResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Pair.RemoteDomain)
				require.NotNil(t, resp.Pair.RemoteToken)
				require.NotNil(t, resp.Pair.LocalToken)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Pair),
				)
			}
		})
	}
}

func TestListTokenPairs(t *testing.T) {
	net, objs := networkWithTokenPairObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTokenPairs(), args)
			require.NoError(t, err)
			var resp types.QueryAllTokenPairsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.TokenPairs), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.TokenPairs),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTokenPairs(), args)
			require.NoError(t, err)
			var resp types.QueryAllTokenPairsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.TokenPairs), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.TokenPairs),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListTokenPairs(), args)
		require.NoError(t, err)
		var resp types.QueryAllTokenPairsResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.TokenPairs),
		)
	})
}
