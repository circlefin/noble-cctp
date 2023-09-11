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

	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

/*
 * Happy path
 * Authority not set
 * Invalid authority
 */

func TestSetMaxBurnAmountPerMessageHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgSetMaxBurnAmountPerMessage{
		From:       tokenController,
		LocalToken: "uusdc",
		Amount:     math.NewInt(123),
	}

	_, err := server.SetMaxBurnAmountPerMessage(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetPerMessageBurnLimit(ctx, message.LocalToken)
	require.True(t, found)
	require.Equal(t, message.LocalToken, actual.Denom)
	require.Equal(t, message.Amount, actual.Amount)
}

func TestSetMaxBurnAmountPerMessageAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgSetMaxBurnAmountPerMessage{
		From:       sample.AccAddress(),
		LocalToken: "uusdc",
		Amount:     math.NewInt(123),
	}

	require.PanicsWithValue(t, "cctp token controller not found in state", func() {
		_, _ = server.SetMaxBurnAmountPerMessage(sdk.WrapSDKContext(ctx), &message)
	})
}

func TestSetMaxBurnAmountPerMessageInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgSetMaxBurnAmountPerMessage{
		From:       "not authority",
		LocalToken: "uusdc",
		Amount:     math.NewInt(123),
	}

	_, err := server.SetMaxBurnAmountPerMessage(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot set the max burn amount per message")
}
