package keeper_test

import (
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/sample"
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
 * Happy path
 * Authority not set
 * Invalid authority
 * Existing token pair found
 */
func TestLinkTokenPairHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: sample.AccAddress()}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgLinkTokenPair{
		From:         authority.Address,
		RemoteDomain: 1,
		RemoteToken:  "0xabcd",
		LocalToken:   "uusdc",
	}

	_, err := server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	actual, found := testkeeper.GetTokenPair(ctx, message.RemoteDomain, message.RemoteToken)
	require.True(t, found)
	require.Equal(t, message.RemoteToken, actual.RemoteToken)
	require.Equal(t, message.RemoteDomain, actual.RemoteDomain)
	require.Equal(t, message.LocalToken, actual.LocalToken)
}

func TestLinkTokenPairAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgLinkTokenPair{
		From:         sample.AccAddress(),
		RemoteDomain: 1,
		RemoteToken:  "0xABCD",
		LocalToken:   "uusdc",
	}

	_, err := server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority is not set")
}

func TestLinkTokenPairInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: "authority"}
	testkeeper.SetAuthority(ctx, authority)

	message := types.MsgLinkTokenPair{
		From:         "not authority",
		RemoteDomain: 1,
		RemoteToken:  "0xABCD",
		LocalToken:   "uusdc",
	}

	_, err := server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot link token pairs")
}

func TestLinkTokenPairExistingTokenPairFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	authority := types.Authority{Address: "authority"}
	testkeeper.SetAuthority(ctx, authority)

	existingTokenPair := types.TokenPair{
		RemoteDomain: 1,
		RemoteToken:  "0xABCD",
		LocalToken:   "uusdc",
	}

	testkeeper.SetTokenPair(ctx, existingTokenPair)

	message := types.MsgLinkTokenPair{
		From:         authority.Address,
		RemoteDomain: existingTokenPair.RemoteDomain,
		RemoteToken:  existingTokenPair.RemoteToken,
		LocalToken:   existingTokenPair.LocalToken,
	}

	_, err := server.LinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrTokenPairAlreadyFound, err)
	require.Contains(t, err.Error(), "local token for this remote domain + remote token mapping already exists in store")
}
