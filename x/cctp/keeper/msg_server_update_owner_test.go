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
 * Invalid new owner
 */

func TestUpdateAuthorityHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	newOwner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgUpdateOwner{
		From:     owner,
		NewOwner: newOwner,
	}
	_, err := server.UpdateOwner(ctx, &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetPendingOwner(ctx)
	require.True(t, found)
	require.Equal(t, message.NewOwner, actual)
}

func TestUpdateAuthorityAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUpdateOwner{
		From:     sample.AccAddress(),
		NewOwner: "new address",
	}

	require.PanicsWithValue(t, "cctp owner not found in state", func() {
		_, _ = server.UpdateOwner(ctx, &message)
	})
}

func TestUpdateAuthorityInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgUpdateOwner{
		From:     "not the authority",
		NewOwner: "new address",
	}
	_, err := server.UpdateOwner(ctx, &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the authority")
}

func TestUpdateAuthorityInvalidNewOwner(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgUpdateOwner{
		From:     owner,
		NewOwner: "invalid address",
	}
	_, err := server.UpdateOwner(ctx, &message)
	require.ErrorIs(t, err, types.ErrInvalidAddress)
	require.Contains(t, err.Error(), "invalid new owner address")
}
