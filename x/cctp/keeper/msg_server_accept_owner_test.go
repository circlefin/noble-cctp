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
 * Owner not set
 * Pending owner not set
 * Invalid Pending owner
 */

func TestAcceptOwnerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	pendingOwner := sample.AccAddress()
	testkeeper.SetPendingOwner(ctx, pendingOwner)

	message := types.MsgAcceptOwner{
		From: pendingOwner,
	}

	_, err := server.AcceptOwner(ctx, &message)
	require.Nil(t, err)

	newOwner := testkeeper.GetOwner(ctx)
	require.Equal(t, pendingOwner, newOwner)
}

func TestAcceptOwnerOwnerNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	pendingOwner := sample.AccAddress()
	testkeeper.SetPendingOwner(ctx, pendingOwner)

	message := types.MsgAcceptOwner{
		From: pendingOwner,
	}

	require.Panicsf(t, func() {
		_, _ = server.AcceptOwner(ctx, &message)
	}, "cctp owner not found in state")
}

func TestAcceptOwnerPendingOwnerNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgAcceptOwner{
		From: sample.AccAddress(),
	}

	_, err := server.AcceptOwner(ctx, &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "pending owner is not set")
}

func TestAcceptOwnerInvalidPendingOwner(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	pendingOwner := sample.AccAddress()
	testkeeper.SetPendingOwner(ctx, pendingOwner)

	message := types.MsgAcceptOwner{
		From: sample.AccAddress(),
	}

	_, err := server.AcceptOwner(ctx, &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "you are not the pending owner")
}
