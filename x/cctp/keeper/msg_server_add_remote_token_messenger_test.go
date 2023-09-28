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

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var tokenMessenger = make([]byte, 32)

func init() {
	tokenMessenger = common.FromHex("0x000000000000000000000000d0c3da58f55358142b8d3e06c1c30c5c6114efe8")
}

/*
* Happy path
* Authority not set
* Invalid authority
* Remote token messenger already found
 */

func TestAddRemoteTokenMessengerHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgAddRemoteTokenMessenger{
		From:     owner,
		DomainId: 0,
		Address:  tokenMessenger,
	}

	_, err := server.AddRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetRemoteTokenMessenger(ctx, message.DomainId)
	require.True(t, found)

	require.Equal(t, message.DomainId, actual.DomainId)
	require.Equal(t, message.Address, actual.Address)
}

func TestAddRemoteTokenMessengerAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgAddRemoteTokenMessenger{
		From:     sample.AccAddress(),
		DomainId: 0,
		Address:  tokenMessenger,
	}

	require.PanicsWithValue(t, "cctp owner not found in state", func() {
		_, _ = server.AddRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	})
}

func TestAddRemoteTokenMessengerInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgAddRemoteTokenMessenger{
		From:     "not the authority address",
		DomainId: 0,
		Address:  tokenMessenger,
	}

	_, err := server.AddRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot add remote token messengers")
}

func TestAddRemoteTokenMessengerTokenMessengerAlreadyFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	existingRemoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: 0,
		Address:  tokenMessenger,
	}
	testkeeper.SetRemoteTokenMessenger(ctx, existingRemoteTokenMessenger)

	message := types.MsgAddRemoteTokenMessenger{
		From:     owner,
		DomainId: existingRemoteTokenMessenger.DomainId,
		Address:  tokenMessenger,
	}

	_, err := server.AddRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrRemoteTokenMessengerAlreadyFound, err)
	require.Contains(t, err.Error(), "a remote token messenger for this domain already exists")
}

func TestAddRemoteTokenMessengerInvalidAddress(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	owner := sample.AccAddress()
	testkeeper.SetOwner(ctx, owner)

	message := types.MsgAddRemoteTokenMessenger{
		From:     owner,
		DomainId: 0,
		Address:  common.FromHex("0xD0C3da58f55358142b8d3e06C1C30c5C6114EFE8"),
	}

	_, err := server.AddRemoteTokenMessenger(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
}
