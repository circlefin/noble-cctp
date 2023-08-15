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

func createNAttesters(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Attester {
	items := make([]types.Attester, n)
	for i := range items {
		items[i].Attester = "Attester" + strconv.Itoa(i)
		keeper.SetAttester(ctx, items[i])
	}
	return items
}

func TestAttestersGet(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNAttesters(cctpKeeper, ctx, 10)
	for _, item := range items {
		attester, found := cctpKeeper.GetAttester(ctx,
			item.Attester,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&attester),
		)
	}
}

func TestAttestersRemove(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNAttesters(cctpKeeper, ctx, 10)
	for _, item := range items {
		cctpKeeper.DeleteAttester(ctx, item.Attester)
		_, found := cctpKeeper.GetAttester(ctx, item.Attester)
		require.False(t, found)
	}
}

func TestAttestersGetAll(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNAttesters(cctpKeeper, ctx, 10)
	denom := make([]types.Attester, len(items))
	copy(denom, items)
	require.ElementsMatch(t,
		nullify.Fill(denom),
		nullify.Fill(cctpKeeper.GetAllAttesters(ctx)),
	)
}
