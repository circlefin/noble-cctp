package cctp

import (
	"math/rand"

	"github.com/circlefin/noble-cctp/x/cctp/simulation"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simTypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

var _ module.AppModuleSimulation = AppModule{}

func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
	simulation.GenerateGenesisState(input)
}

func (am AppModule) ProposalContents(_ module.SimulationState) []simTypes.WeightedProposalContent {
	// We don't have any governance proposals in the CCTP module.
	return nil
}

func (am AppModule) RandomizedParams(_ *rand.Rand) []simTypes.ParamChange {
	// We don't have any parameters in the CCTP module.
	return nil
}

func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

func (am AppModule) WeightedOperations(simState module.SimulationState) []simTypes.WeightedOperation {
	return simulation.WeightedOperations(simState.Cdc, am.accountKeeper, am.bankKeeper, am.keeper)
}
