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
 * Invalid new token controller address
 */

func TestUpdateTokenControllerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	newTokenController := sample.AccAddress()

	message := types.MsgUpdateTokenController{
		From:               owner,
		NewTokenController: newTokenController,
	}

	_, err := server.UpdateTokenController(ctx, &message)
	require.Nil(t, err)

	actual := testkeeper.GetTokenController(ctx)
	require.Equal(t, newTokenController, actual)
}

func TestUpdateTokenControllerAuthorityIsNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgUpdateTokenController{
		From:               "not the authority",
		NewTokenController: sample.AccAddress(),
	}
	require.Panicsf(t, func() {
		_, _ = server.UpdateTokenController(ctx, &message)
	}, "cctp owner not found in state")
}

func TestUpdateTokenControllerInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	newTokenController := sample.AccAddress()

	message := types.MsgUpdateTokenController{
		From:               sample.AccAddress(),
		NewTokenController: newTokenController,
	}

	_, err := server.UpdateTokenController(ctx, &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the authority")
}

func TestUpdateTokenControllerInvalidNewTokenController(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	newTokenController := "invalid new token controller"

	message := types.MsgUpdateTokenController{
		From:               owner,
		NewTokenController: newTokenController,
	}

	_, err := server.UpdateTokenController(ctx, &message)
	require.ErrorIs(t, err, types.ErrInvalidAddress)
	require.Contains(t, err.Error(), "invalid new token controller address")
}
