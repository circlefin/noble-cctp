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
 * Attester already found
 */
func TestEnableAttesterHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgEnableAttester{
		From:     attesterManager,
		Attester: "attester",
	}

	_, err := server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetAttester(ctx, message.Attester)
	require.True(t, found)
	require.Equal(t, message.Attester, actual.Attester)
}

func TestEnableAttesterAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgEnableAttester{
		From:     sample.AccAddress(),
		Attester: "attester",
	}

	require.PanicsWithValue(t, "cctp attester manager not found in state", func() {
		_, _ = server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	})
}

func TestEnableAttesterInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgEnableAttester{
		From:     sample.AccAddress(),
		Attester: "attester",
	}

	_, err := server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot enable attesters")
}

func TestEnableAttesterAttesterAlreadyFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existingAttester := types.Attester{Attester: "attester"}
	testkeeper.SetAttester(ctx, existingAttester)

	message := types.MsgEnableAttester{
		From:     attesterManager,
		Attester: "attester",
	}

	_, err := server.EnableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAttesterAlreadyFound, err)
	require.Contains(t, err.Error(), "this attester already exists in the store")
}
