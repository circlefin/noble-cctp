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
package keeper_test

import (
	"strconv"
	"testing"

	"github.com/circlefin/noble-cctp/x/cctp/keeper"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNTokenPairs(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.TokenPair {
	items := make([]types.TokenPair, n)
	for i := range items {
		items[i].RemoteDomain = uint32(i)
		items[i].RemoteToken = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, byte(i)}
		items[i].LocalToken = "token" + strconv.Itoa(i)

		keeper.SetTokenPair(ctx, items[i])
	}
	return items
}

func TestTokenPairsGet(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNTokenPairs(cctpKeeper, ctx, 10)
	for _, item := range items {
		tokenPair, found := cctpKeeper.GetTokenPair(ctx,
			item.RemoteDomain,
			item.RemoteToken,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&tokenPair),
		)
	}
}

func TestTokenPairsRemove(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNTokenPairs(cctpKeeper, ctx, 10)
	for _, item := range items {
		cctpKeeper.DeleteTokenPair(
			ctx,
			item.RemoteDomain,
			item.RemoteToken,
		)
		_, found := cctpKeeper.GetTokenPair(
			ctx,
			item.RemoteDomain,
			item.RemoteToken,
		)
		require.False(t, found)
	}
}

func TestTokenPairsGetAll(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNTokenPairs(cctpKeeper, ctx, 10)
	denom := make([]types.TokenPair, len(items))
	copy(denom, items)

	require.ElementsMatch(t,
		nullify.Fill(denom),
		nullify.Fill(cctpKeeper.GetAllTokenPairs(ctx)),
	)
}
