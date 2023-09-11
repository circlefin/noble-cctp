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
 * Attester not found
 * Fails when only 1 attester is left
 * Fails when signature threshold not found
 * Fails when signature threshold is too low
 */
func TestDisableAttesterHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "attester",
	}
	existing2 := types.Attester{
		Attester: "attester2",
	}
	existing3 := types.Attester{
		Attester: "attester3",
	}
	testkeeper.SetAttester(ctx, existing)
	testkeeper.SetAttester(ctx, existing2)
	testkeeper.SetAttester(ctx, existing3)

	sig := types.SignatureThreshold{Amount: 2}
	testkeeper.SetSignatureThreshold(ctx, sig)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "attester",
	}

	_, err := server.DisableAttester(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	_, found := testkeeper.GetAttester(ctx, message.Attester)
	require.False(t, found)
}

func TestDisableAttesterAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	existing := types.Attester{
		Attester: "attester",
	}
	testkeeper.SetAttester(ctx, existing)

	message := types.MsgDisableAttester{
		From:     sample.AccAddress(),
		Attester: "attester",
	}

	require.PanicsWithValue(t, "cctp attester manager not found in state", func() {
		_, _ = server.DisableAttester(sdk.WrapSDKContext(ctx), &message)
	})
}

func TestDisableAttesterInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "attester",
	}
	testkeeper.SetAttester(ctx, existing)

	message := types.MsgDisableAttester{
		From:     sample.AccAddress(),
		Attester: "attester",
	}

	_, err := server.DisableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot disable attesters")
}

func TestDisableAttesterAttesterNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "attester",
	}

	_, err := server.DisableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "attester not found")
}

func TestDisableAttesterFailsWhenOnly1AttesterIsLeft(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "attester",
	}
	testkeeper.SetAttester(ctx, existing)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "attester",
	}

	_, err := server.DisableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "cannot disable the last attester")

	_, found := testkeeper.GetAttester(ctx, message.Attester)
	require.True(t, found)
}

func TestDisableAttesterFailsWhenSignatureThresholdNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "attester",
	}
	existing2 := types.Attester{
		Attester: "attester2",
	}
	testkeeper.SetAttester(ctx, existing)
	testkeeper.SetAttester(ctx, existing2)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "attester",
	}

	_, err := server.DisableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "signature threshold not set")
}

func TestDisableAttesterFailsWhenSignatureThresholdIsTooLow(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing1 := types.Attester{
		Attester: "attester1",
	}
	existing2 := types.Attester{
		Attester: "attester2",
	}
	testkeeper.SetAttester(ctx, existing1)
	testkeeper.SetAttester(ctx, existing2)

	sig := types.SignatureThreshold{Amount: 2}
	testkeeper.SetSignatureThreshold(ctx, sig)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "attester1",
	}

	_, err := server.DisableAttester(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "signature threshold is too low")
}
