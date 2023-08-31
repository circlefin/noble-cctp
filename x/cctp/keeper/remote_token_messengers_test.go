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
	"testing"

	"github.com/circlefin/noble-cctp/x/cctp/keeper"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func createNRemoteTokenMessengers(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.RemoteTokenMessenger {
	items := make([]types.RemoteTokenMessenger, n)
	for i := range items {
		items[i].DomainId = uint32(i)
		items[i].Address = make([]byte, 32)

		keeper.SetRemoteTokenMessenger(ctx, items[i])
	}
	return items
}

func TestRemoteTokenMessengersGet(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNRemoteTokenMessengers(cctpKeeper, ctx, 10)
	for _, item := range items {
		tokenMessenger, found := cctpKeeper.GetRemoteTokenMessenger(ctx, item.DomainId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&tokenMessenger),
		)
	}
}

func TestRemoteTokenMessengersRemove(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNRemoteTokenMessengers(cctpKeeper, ctx, 10)
	for _, item := range items {
		cctpKeeper.DeleteRemoteTokenMessenger(ctx, item.DomainId)
		_, found := cctpKeeper.GetRemoteTokenMessenger(ctx, item.DomainId)
		require.False(t, found)
	}
}

func TestRemoteTokenMessengersGetAll(t *testing.T) {
	cctpKeeper, ctx := keepertest.CctpKeeper(t)
	items := createNRemoteTokenMessengers(cctpKeeper, ctx, 10)
	denom := make([]types.RemoteTokenMessenger, len(items))
	copy(denom, items)

	require.ElementsMatch(t,
		nullify.Fill(denom),
		nullify.Fill(cctpKeeper.GetRemoteTokenMessengers(ctx)),
	)
}
