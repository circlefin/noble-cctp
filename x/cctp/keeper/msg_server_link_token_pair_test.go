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
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var token = make([]byte, 32)

func init() {
	token = common.FromHex("0x00000000000000000000000007865c6e87b9f70255377e024ace6630c1eaa37f")
}

/*
 * Happy path
 * Authority not set
 * Invalid authority
 * Existing token pair found
 */
func TestLinkTokenPairHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgLinkTokenPair{
		From:         tokenController,
		RemoteDomain: 0,
		RemoteToken:  token,
		LocalToken:   "uusdc",
	}

	_, err := server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.NoError(t, err)

	actual, found := testkeeper.GetTokenPair(ctx, message.RemoteDomain, message.RemoteToken)
	require.True(t, found)
	require.Equal(t, message.RemoteDomain, actual.RemoteDomain)
	require.Equal(t, message.RemoteToken, actual.RemoteToken)
	require.Equal(t, message.LocalToken, actual.LocalToken)
}

func TestLinkTokenPairAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgLinkTokenPair{
		From:         sample.AccAddress(),
		RemoteDomain: 0,
		RemoteToken:  token,
		LocalToken:   "uusdc",
	}

	require.PanicsWithValue(t, "cctp token controller not found in state", func() {
		_, _ = server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	})
}

func TestLinkTokenPairInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgLinkTokenPair{
		From:         "not authority",
		RemoteDomain: 0,
		RemoteToken:  token,
		LocalToken:   "uusdc",
	}

	_, err := server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
}

func TestLinkTokenPairExistingTokenPairFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	existingTokenPair := types.TokenPair{
		RemoteDomain: 0,
		RemoteToken:  token,
		LocalToken:   "uusdc",
	}

	testkeeper.SetTokenPair(ctx, existingTokenPair)

	message := types.MsgLinkTokenPair{
		From:         tokenController,
		RemoteDomain: existingTokenPair.RemoteDomain,
		RemoteToken:  existingTokenPair.RemoteToken,
		LocalToken:   existingTokenPair.LocalToken,
	}

	_, err := server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrTokenPairAlreadyFound, err)
}
