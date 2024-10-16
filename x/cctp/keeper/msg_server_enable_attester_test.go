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
 * Attester already found
 */
func TestEnableAttesterHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgEnableAttester{
		From:     attesterManager,
		Attester: "1234",
	}

	_, err := server.EnableAttester(ctx, &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetAttester(ctx, message.Attester)
	require.True(t, found)
	require.Equal(t, message.Attester, actual.Attester)
}

func TestEnableAttesterAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgEnableAttester{
		From:     sample.AccAddress(),
		Attester: "1234",
	}

	require.PanicsWithValue(t, "cctp attester manager not found in state", func() {
		_, _ = server.EnableAttester(ctx, &message)
	})
}

func TestEnableAttesterInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgEnableAttester{
		From:     sample.AccAddress(),
		Attester: "1234",
	}

	_, err := server.EnableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot enable attesters")
}

func TestEnableAttesterInvalidAttester(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	message := types.MsgEnableAttester{
		From:     attesterManager,
		Attester: "invalid attester",
	}

	_, err := server.EnableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrInvalidAddress, err)
	require.Contains(t, err.Error(), "invalid attester")
}

func TestEnableAttesterAttesterAlreadyFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	attesterManager := sample.AccAddress()
	testkeeper.SetAttesterManager(ctx, attesterManager)

	existingAttester := types.Attester{Attester: "1234"}
	testkeeper.SetAttester(ctx, existingAttester)

	message := types.MsgEnableAttester{
		From:     attesterManager,
		Attester: existingAttester.Attester,
	}

	_, err := server.EnableAttester(ctx, &message)
	require.ErrorIs(t, types.ErrAttesterAlreadyFound, err)
	require.Contains(t, err.Error(), "this attester already exists in the store")
}
