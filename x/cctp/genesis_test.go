package cctp_test

import (
	"testing"

	"cosmossdk.io/math"
	keepertest "github.com/circlefin/noble-cctp/testutil/keeper"
	"github.com/circlefin/noble-cctp/testutil/nullify"
	"github.com/circlefin/noble-cctp/x/cctp"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenesisHappyPath(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Authority: &types.Authority{
			Address: "123",
		},
		AttesterList: []types.Attester{
			{
				Attester: "0",
			},
			{
				Attester: "1",
			},
		},
		PerMessageBurnLimitList: []types.PerMessageBurnLimit{
			{
				Denom:  "uusdc",
				Amount: math.NewInt(int64(1)),
			},
			{
				Denom:  "euroc",
				Amount: math.NewInt(int64(2)),
			},
		},
		BurningAndMintingPaused: &types.BurningAndMintingPaused{
			Paused: true,
		},
		SendingAndReceivingMessagesPaused: &types.SendingAndReceivingMessagesPaused{
			Paused: false,
		},
		MaxMessageBodySize: &types.MaxMessageBodySize{
			Amount: 12,
		},
		NextAvailableNonce: &types.Nonce{
			Nonce: 34,
		},
		SignatureThreshold: &types.SignatureThreshold{
			Amount: 2,
		},
		TokenPairList: []types.TokenPair{
			{
				RemoteDomain: uint32(0),
				RemoteToken:  []byte("1"),
				LocalToken:   "uusdc",
			},
			{
				RemoteDomain: uint32(1),
				RemoteToken:  []byte("2"),
				LocalToken:   "uusdc",
			},
		},
		UsedNoncesList: []types.Nonce{
			{
				SourceDomain: uint32(1),
				Nonce:        uint64(1234),
			},
			{
				SourceDomain: uint32(2),
				Nonce:        uint64(5678),
			},
		},
		TokenMessengerList: []types.TokenMessenger{
			{
				DomainId: uint32(1),
				Address:  "1234",
			},
			{
				DomainId: uint32(2),
				Address:  "56789",
			},
		},
	}

	k, ctx := keepertest.CctpKeeper(t)
	cctp.InitGenesis(ctx, k, genesisState)
	got := cctp.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.Authority, got.Authority)
	require.ElementsMatch(t, genesisState.AttesterList, got.AttesterList)
	require.ElementsMatch(t, genesisState.PerMessageBurnLimitList, got.PerMessageBurnLimitList)
	require.Equal(t, genesisState.BurningAndMintingPaused, got.BurningAndMintingPaused)
	require.Equal(t, genesisState.SendingAndReceivingMessagesPaused, got.SendingAndReceivingMessagesPaused)
	require.Equal(t, genesisState.MaxMessageBodySize, got.MaxMessageBodySize)
	require.Equal(t, genesisState.NextAvailableNonce, got.NextAvailableNonce)
	require.Equal(t, genesisState.SignatureThreshold, got.SignatureThreshold)
	require.ElementsMatch(t, genesisState.TokenPairList, got.TokenPairList)
	require.ElementsMatch(t, genesisState.UsedNoncesList, got.UsedNoncesList)
	require.ElementsMatch(t, genesisState.TokenMessengerList, got.TokenMessengerList)
}

func TestGenesisPanicsWithNoAuthority(t *testing.T) {
	genesisState := types.GenesisState{}

	k, ctx := keepertest.CctpKeeper(t)

	assert.Panics(t, func() {
		cctp.InitGenesis(ctx, k, genesisState)
	})
}

func TestGenesisBurningAndMintingPausedDefault(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Authority: &types.Authority{
			Address: "123",
		},
	}
	k, ctx := keepertest.CctpKeeper(t)

	cctp.InitGenesis(ctx, k, genesisState)
	got := cctp.ExportGenesis(ctx, k)

	require.Equal(t, true, got.BurningAndMintingPaused.Paused)
}

func TestGenesisSendingAndReceivingMessagesPausedDefault(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Authority: &types.Authority{
			Address: "123",
		},
		BurningAndMintingPaused: &types.BurningAndMintingPaused{Paused: true},
	}
	k, ctx := keepertest.CctpKeeper(t)

	cctp.InitGenesis(ctx, k, genesisState)
	got := cctp.ExportGenesis(ctx, k)

	require.Equal(t, true, got.SendingAndReceivingMessagesPaused.Paused)
}

func TestGenesisMaxMessageBodySizeIsDefault(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Authority: &types.Authority{
			Address: "123",
		},
		BurningAndMintingPaused:           &types.BurningAndMintingPaused{Paused: true},
		SendingAndReceivingMessagesPaused: &types.SendingAndReceivingMessagesPaused{Paused: true},
	}
	k, ctx := keepertest.CctpKeeper(t)

	cctp.InitGenesis(ctx, k, genesisState)
	got := cctp.ExportGenesis(ctx, k)

	require.Equal(t, uint64(8000), got.MaxMessageBodySize.Amount)
}

func TestGenesisNextAvailableNonceDefault(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Authority: &types.Authority{
			Address: "123",
		},
		BurningAndMintingPaused:           &types.BurningAndMintingPaused{Paused: true},
		SendingAndReceivingMessagesPaused: &types.SendingAndReceivingMessagesPaused{Paused: true},
	}
	k, ctx := keepertest.CctpKeeper(t)

	cctp.InitGenesis(ctx, k, genesisState)
	got := cctp.ExportGenesis(ctx, k)

	require.Equal(t, uint64(0), got.NextAvailableNonce.Nonce)
}

func TestGenesisSignatureThresholdPanicsWhenZero(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Authority: &types.Authority{
			Address: "123",
		},
		BurningAndMintingPaused:           &types.BurningAndMintingPaused{Paused: true},
		SendingAndReceivingMessagesPaused: &types.SendingAndReceivingMessagesPaused{Paused: true},
		SignatureThreshold:                &types.SignatureThreshold{Amount: uint32(0)},
	}
	k, ctx := keepertest.CctpKeeper(t)

	assert.Panics(t, func() {
		cctp.InitGenesis(ctx, k, genesisState)
	})
}

func TestGenesisSignatureThresholdDefault(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Authority: &types.Authority{
			Address: "123",
		},
		BurningAndMintingPaused:           &types.BurningAndMintingPaused{Paused: true},
		SendingAndReceivingMessagesPaused: &types.SendingAndReceivingMessagesPaused{Paused: true},
	}
	k, ctx := keepertest.CctpKeeper(t)

	cctp.InitGenesis(ctx, k, genesisState)
	got := cctp.ExportGenesis(ctx, k)

	require.Equal(t, uint32(1), got.SignatureThreshold.Amount)
}
