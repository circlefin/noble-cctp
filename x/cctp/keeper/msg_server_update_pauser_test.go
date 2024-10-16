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
 * Invalid new pauser address
 */

func TestUpdatePauserHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	newPauser := sample.AccAddress()

	message := types.MsgUpdatePauser{
		From:      owner,
		NewPauser: newPauser,
	}

	_, err := server.UpdatePauser(ctx, &message)
	require.Nil(t, err)

	actual := testkeeper.GetPauser(ctx)
	require.Equal(t, newPauser, actual)
}

func TestUpdatePauserAuthorityIsNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	message := types.MsgUpdatePauser{
		From:      "not the authority",
		NewPauser: sample.AccAddress(),
	}
	require.Panicsf(t, func() {
		_, _ = server.UpdatePauser(ctx, &message)
	}, "cctp owner not found in state")
}

func TestUpdatePauserInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	newPauser := sample.AccAddress()

	message := types.MsgUpdatePauser{
		From:      sample.AccAddress(),
		NewPauser: newPauser,
	}

	_, err := server.UpdatePauser(ctx, &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot update the pauser")
}

func TestUpdatePauserInvalidNewPauser(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper()
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)
	pauser := sample.AccAddress()
	testkeeper.SetPauser(ctx, pauser)

	newPauser := "invalid new pauser"

	message := types.MsgUpdatePauser{
		From:      owner,
		NewPauser: newPauser,
	}

	_, err := server.UpdatePauser(ctx, &message)
	require.ErrorIs(t, err, types.ErrInvalidAddress)
	require.Contains(t, err.Error(), "invalid new pauser address")
}
