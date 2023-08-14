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
 * Token pair not found
 */
func TestUnlinkTokenPairHappyPath(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	tokenPair := types.TokenPair{
		RemoteDomain: 1,
		RemoteToken:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xAB, 0xCD},
		LocalToken:   "uusdc",
	}
	testkeeper.SetTokenPair(ctx, tokenPair)

	message := types.MsgUnlinkTokenPair{
		From:         tokenController,
		RemoteDomain: 1,
		RemoteToken:  "0xABCD",
		LocalToken:   "uusdc",
	}

	_, err := server.UnlinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.Nil(t, err)

	_, found := testkeeper.GetTokenPairHex(ctx, message.RemoteDomain, message.RemoteToken)
	require.False(t, found)
}

func TestUnlinkTokenPairAuthorityNotSet(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	message := types.MsgUnlinkTokenPair{
		From:         sample.AccAddress(),
		RemoteDomain: 1,
		RemoteToken:  "0xABCD",
		LocalToken:   "uusdc",
	}

	_, err := server.UnlinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrAuthorityNotSet, err)
	require.Contains(t, err.Error(), "authority not set")
}

func TestUnlinkTokenPairInvalidAuthority(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgUnlinkTokenPair{
		From:         "not the authority",
		RemoteDomain: 1,
		RemoteToken:  "0xABCD",
		LocalToken:   "uusdc",
	}

	_, err := server.UnlinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrUnauthorized, err)
	require.Contains(t, err.Error(), "this message sender cannot unlink token pairs")
}

func TestUnlinkTokenPairTokenPairNotFound(t *testing.T) {
	testkeeper, ctx := keepertest.CctpKeeper(t)
	server := keeper.NewMsgServerImpl(testkeeper)

	tokenController := sample.AccAddress()
	testkeeper.SetTokenController(ctx, tokenController)

	message := types.MsgUnlinkTokenPair{
		From:         tokenController,
		RemoteDomain: 1,
		RemoteToken:  "0xABCD",
		LocalToken:   "uusdc",
	}

	_, err := server.UnlinkTokenPair(sdk.WrapSDKContext(ctx), &message)
	require.ErrorIs(t, types.ErrTokenPairNotFound, err)
	require.Contains(t, err.Error(), "token pair doesn't exist in store")
}
