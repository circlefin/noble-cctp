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

	"github.com/stretchr/testify/require"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func TestNextAvailableNonce(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)

	_, found := keeper.GetNextAvailableNonce(ctx)
	require.False(t, found)

	savedNonce := types.Nonce{Nonce: 21}
	keeper.SetNextAvailableNonce(ctx, savedNonce)

	next, found := keeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t,
		savedNonce,
		nullify.Fill(&next),
	)

	newSavedNonce := types.Nonce{Nonce: 22}

	keeper.SetNextAvailableNonce(ctx, newSavedNonce)

	next, found = keeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t,
		newSavedNonce,
		nullify.Fill(&next),
	)
}

func TestNextAvailableNonceReserveAndIncrement(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)

	savedNonce := types.Nonce{Nonce: 21}
	keeper.SetNextAvailableNonce(ctx, savedNonce)

	prev, found := keeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t,
		savedNonce,
		nullify.Fill(&prev),
	)

	// method returns the nonce being reserved
	nextFromMethod := keeper.ReserveAndIncrementNonce(ctx)
	require.Equal(t,
		types.Nonce{
			Nonce: prev.Nonce,
		},
		nullify.Fill(&nextFromMethod),
	)

	// retrieving the nonce should get reserved nonce + 1
	next, found := keeper.GetNextAvailableNonce(ctx)
	require.True(t, found)
	require.Equal(t,
		types.Nonce{
			Nonce: prev.Nonce + 1,
		},
		nullify.Fill(&next),
	)
}
