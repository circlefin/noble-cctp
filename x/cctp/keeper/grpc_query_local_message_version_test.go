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

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestLocalMessageVersionQuery(t *testing.T) {
	keeper, ctx := keepertest.CctpKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)

	res, err := keeper.LocalMessageVersion(goCtx, &types.QueryLocalMessageVersionRequest{})

	require.NoError(t, err)
	require.Equal(t, res, &types.QueryLocalMessageVersionResponse{Version: types.NobleMessageVersion})
}
