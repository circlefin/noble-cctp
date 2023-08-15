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
	"github.com/circlefin/noble-cctp/x/cctp/keeper"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	_ "github.com/cosmos/cosmos-sdk/types/errors" // sdkerrors

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genState types.GenesisState) {
	k.SetOwner(ctx, genState.Owner)
	k.SetAttesterManager(ctx, genState.AttesterManager)
	k.SetPauser(ctx, genState.Pauser)
	k.SetTokenController(ctx, genState.TokenController)

	for _, elem := range genState.AttesterList {
		k.SetAttester(ctx, elem)
	}

	for _, elem := range genState.PerMessageBurnLimitList {
		k.SetPerMessageBurnLimit(ctx, elem)
	}

	if genState.BurningAndMintingPaused != nil {
		k.SetBurningAndMintingPaused(ctx, *genState.BurningAndMintingPaused)
	} else {
		k.SetBurningAndMintingPaused(ctx, types.BurningAndMintingPaused{Paused: true})
	}

	if genState.SendingAndReceivingMessagesPaused != nil {
		k.SetSendingAndReceivingMessagesPaused(ctx, *genState.SendingAndReceivingMessagesPaused)
	} else {
		k.SetSendingAndReceivingMessagesPaused(ctx, types.SendingAndReceivingMessagesPaused{Paused: true})
	}

	if genState.MaxMessageBodySize != nil {
		k.SetMaxMessageBodySize(ctx, *genState.MaxMessageBodySize)
	} else {
		k.SetMaxMessageBodySize(ctx, types.MaxMessageBodySize{Amount: 8000})
	}

	if genState.NextAvailableNonce != nil {
		k.SetNextAvailableNonce(ctx, *genState.NextAvailableNonce)
	} else {
		k.SetNextAvailableNonce(ctx, types.Nonce{Nonce: 0})
	}

	if genState.SignatureThreshold != nil {
		if genState.SignatureThreshold.Amount == 0 {
			panic("Signature threshold must not be 0")
		}
		k.SetSignatureThreshold(ctx, *genState.SignatureThreshold)
	} else {
		k.SetSignatureThreshold(ctx, types.SignatureThreshold{Amount: 1})
	}

	for _, elem := range genState.TokenPairList {
		k.SetTokenPair(ctx, elem)
	}

	for _, elem := range genState.UsedNoncesList {
		k.SetUsedNonce(ctx, elem)
	}

	for _, elem := range genState.TokenMessengerList {
		k.SetRemoteTokenMessenger(ctx, elem)
	}
}

// ExportGenesis returns the module's exported GenesisState
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Owner = k.GetOwner(ctx)
	genesis.AttesterManager = k.GetAttesterManager(ctx)
	genesis.Pauser = k.GetPauser(ctx)
	genesis.TokenController = k.GetTokenController(ctx)

	genesis.AttesterList = k.GetAllAttesters(ctx)
	genesis.PerMessageBurnLimitList = k.GetAllPerMessageBurnLimits(ctx)

	burningAndMintingPaused, found := k.GetBurningAndMintingPaused(ctx)
	if found {
		genesis.BurningAndMintingPaused = &burningAndMintingPaused
	}

	sendingAndReceivingMessagesPaused, found := k.GetSendingAndReceivingMessagesPaused(ctx)
	if found {
		genesis.SendingAndReceivingMessagesPaused = &sendingAndReceivingMessagesPaused
	}

	maxMessageBodySize, found := k.GetMaxMessageBodySize(ctx)
	if found {
		genesis.MaxMessageBodySize = &maxMessageBodySize
	}

	nextAvailableNonce, found := k.GetNextAvailableNonce(ctx)
	if found {
		genesis.NextAvailableNonce = &nextAvailableNonce
	}

	signatureThreshold, found := k.GetSignatureThreshold(ctx)
	if found {
		genesis.SignatureThreshold = &signatureThreshold
	}

	genesis.TokenPairList = k.GetAllTokenPairs(ctx)
	genesis.UsedNoncesList = k.GetAllUsedNonces(ctx)
	genesis.TokenMessengerList = k.GetRemoteTokenMessengers(ctx)

	return genesis
}
