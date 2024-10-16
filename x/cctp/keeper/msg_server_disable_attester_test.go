// Copyright 2024 Circle Internet Group, Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package keeper_test

import (
	"testing"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/stretchr/testify/require"
)

/*
 * Happy path
 * Authority not set
 * Invalid authority
 * Invalid attester
 * Attester not found
 * Fails when only 1 attester is left
 * Fails when signature threshold not found
 * Fails when signature threshold is too low
 */
func TestDisableAttesterHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "1234",
	}
	existing2 := types.Attester{
		Attester: "5678",
	}
	existing3 := types.Attester{
		Attester: "9012",
	}
	testkeeper.SetAttester(ctx, existing)
	testkeeper.SetAttester(ctx, existing2)
	testkeeper.SetAttester(ctx, existing3)

	sig := types.SignatureThreshold{Amount: 2}
	testkeeper.SetSignatureThreshold(ctx, sig)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "1234",
	}

	_, err := server.DisableAttester(ctx, &message)
	require.Nil(t, err)

	_, found := testkeeper.GetAttester(ctx, message.Attester)
	require.False(t, found)
}

func TestDisableAttesterAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	existing := types.Attester{
		Attester: "1234",
	}
	testkeeper.SetAttester(ctx, existing)

	message := types.MsgDisableAttester{
		From:     sample.AccAddress(),
		Attester: "1234",
	}

	require.PanicsWithValue(t, "cctp attester manager not found in state", func() {
		_, _ = server.DisableAttester(ctx, &message)
	})
}

func TestDisableAttesterInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "1234",
	}
	testkeeper.SetAttester(ctx, existing)

	message := types.MsgDisableAttester{
		From:     sample.AccAddress(),
		Attester: "1234",
	}

	_, err := server.DisableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot disable attesters")
}

func TestDisableAttesterInvalidAttester(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "",
	}

	_, err := server.DisableAttester(ctx, &message)
	require.ErrorIs(t, err, types.ErrInvalidAddress)
	require.Contains(t, err.Error(), "invalid attester")
}

func TestDisableAttesterAttesterNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "1234",
	}

	_, err := server.DisableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "attester not found")
}

func TestDisableAttesterFailsWhenOnly1AttesterIsLeft(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "1234",
	}
	testkeeper.SetAttester(ctx, existing)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "1234",
	}

	_, err := server.DisableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "cannot disable the last attester")

	_, found := testkeeper.GetAttester(ctx, message.Attester)
	require.True(t, found)
}

func TestDisableAttesterFailsWhenSignatureThresholdNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing := types.Attester{
		Attester: "1234",
	}
	existing2 := types.Attester{
		Attester: "5678",
	}
	testkeeper.SetAttester(ctx, existing)
	testkeeper.SetAttester(ctx, existing2)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "1234",
	}

	_, err := server.DisableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "signature threshold not set")
}

func TestDisableAttesterFailsWhenSignatureThresholdIsTooLow(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existing1 := types.Attester{
		Attester: "1234",
	}
	existing2 := types.Attester{
		Attester: "5678",
	}
	testkeeper.SetAttester(ctx, existing1)
	testkeeper.SetAttester(ctx, existing2)

	sig := types.SignatureThreshold{Amount: 2}
	testkeeper.SetSignatureThreshold(ctx, sig)

	message := types.MsgDisableAttester{
		From:     attesterManager,
		Attester: "1234",
	}

	_, err := server.DisableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrDisableAttester, err)
	require.Contains(t, err.Error(), "signature threshold is too low")
}
