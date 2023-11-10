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
