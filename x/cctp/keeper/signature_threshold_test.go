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

func TestSignatureThreshold(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)
	_, found := keeper.GetSignatureThreshold(ctx)
	require.False(t, found)

	SignatureThreshold := types.SignatureThreshold{Amount: 2}
	keeper.SetSignatureThreshold(ctx, SignatureThreshold)

	threshold, found := keeper.GetSignatureThreshold(ctx)
	require.True(t, found)
	require.Equal(t,
		SignatureThreshold,
		nullify.Fill(&threshold),
	)

	newSignatureThreshold := types.SignatureThreshold{Amount: 3}

	keeper.SetSignatureThreshold(ctx, newSignatureThreshold)

	threshold, found = keeper.GetSignatureThreshold(ctx)
	require.True(t, found)
	require.Equal(t,
		newSignatureThreshold,
		nullify.Fill(&threshold),
	)
}
