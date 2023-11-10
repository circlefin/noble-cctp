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
 package simulation

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GenerateGenesisState(simState *module.SimulationState) {
	owner, _ := simulation.RandomAcc(simState.Rand, simState.Accounts)
	attesterManager, _ := simulation.RandomAcc(simState.Rand, simState.Accounts)
	pauser, _ := simulation.RandomAcc(simState.Rand, simState.Accounts)
	tokenController, _ := simulation.RandomAcc(simState.Rand, simState.Accounts)

	genesis := types.GenesisState{
		Owner:           owner.Address.String(),
		AttesterManager: attesterManager.Address.String(),
		Pauser:          pauser.Address.String(),
		TokenController: tokenController.Address.String(),

		BurningAndMintingPaused:           &types.BurningAndMintingPaused{Paused: simState.Rand.Int63n(101) <= 50},
		SendingAndReceivingMessagesPaused: &types.SendingAndReceivingMessagesPaused{Paused: simState.Rand.Int63n(101) > 50},
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&genesis)

	// Give the first account some USDC for testing.
	user := simState.Accounts[0]

	var bankGenesis bankTypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[bankTypes.ModuleName], &bankGenesis)

	for i := 0; i < len(bankGenesis.Balances); i++ {
		balance := bankGenesis.Balances[i]
		if balance.Address == user.Address.String() {
			coin := sdk.NewCoin("uusdc", sdk.NewInt(simState.Rand.Int63()))
			balance.Coins = balance.Coins.Add(coin)

			bankGenesis.Balances[i] = balance
			bankGenesis.Supply = bankGenesis.Supply.Add(coin)
		}
	}

	simState.GenState[bankTypes.ModuleName] = simState.Cdc.MustMarshalJSON(&bankGenesis)
}
