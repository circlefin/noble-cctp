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

package types

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Owner:                             "",
		AttesterManager:                   "",
		Pauser:                            "",
		TokenController:                   "",
		AttesterList:                      []Attester{},
		PerMessageBurnLimitList:           []PerMessageBurnLimit{},
		BurningAndMintingPaused:           &BurningAndMintingPaused{Paused: false},
		SendingAndReceivingMessagesPaused: &SendingAndReceivingMessagesPaused{Paused: false},
		MaxMessageBodySize:                nil,
		NextAvailableNonce:                nil,
		SignatureThreshold:                nil,
		TokenPairList:                     []TokenPair{},
		UsedNoncesList:                    []Nonce{},
		TokenMessengerList:                []RemoteTokenMessenger{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.  Stateful checks are performed in InitGenesis
func (gs GenesisState) Validate() error {
	if gs.Owner != "" {
		if _, err := sdk.AccAddressFromBech32(gs.Owner); err != nil {
			return errors.Wrapf(ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}

	if gs.AttesterManager != "" {
		if _, err := sdk.AccAddressFromBech32(gs.AttesterManager); err != nil {
			return errors.Wrapf(ErrInvalidAddress, "invalid attester manager address (%s)", err)
		}
	}

	if gs.Pauser != "" {
		if _, err := sdk.AccAddressFromBech32(gs.Pauser); err != nil {
			return errors.Wrapf(ErrInvalidAddress, "invalid pauser address (%s)", err)
		}
	}

	if gs.TokenController != "" {
		if _, err := sdk.AccAddressFromBech32(gs.TokenController); err != nil {
			return errors.Wrapf(ErrInvalidAddress, "invalid token controller address (%s)", err)
		}
	}

	// Check for duplicated index in attesters
	attesterIndexMap := make(map[string]struct{})
	for _, elem := range gs.AttesterList {
		index := string(AttesterKey([]byte(elem.Attester)))
		if _, ok := attesterIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for attesters")
		}
		attesterIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in per message burn limit
	perMessageBurnLimitIndexMap := make(map[string]struct{})
	for _, elem := range gs.PerMessageBurnLimitList {
		index := string(PerMessageBurnLimitKey(elem.Denom))
		if _, ok := perMessageBurnLimitIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for per message burn limits")
		}
		perMessageBurnLimitIndexMap[index] = struct{}{}
	}

	if gs.BurningAndMintingPaused == nil {
		return fmt.Errorf("BurningAndMintingPaused cannot be nil")
	}

	if gs.SendingAndReceivingMessagesPaused == nil {
		return fmt.Errorf("SendingAndReceivingMessagesPaused cannot be nil")
	}

	// Check for duplicated index in token pairs
	tokenPairIndexMap := make(map[string]struct{})
	for _, elem := range gs.TokenPairList {
		index := string(TokenPairKey(elem.RemoteDomain, elem.RemoteToken))
		if _, ok := attesterIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for token pairs")
		}
		tokenPairIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in the used nonce list
	usedNonceIndexMap := make(map[string]struct{})
	for _, elem := range gs.UsedNoncesList {
		index := string(UsedNonceKey(elem.Nonce, elem.SourceDomain))
		if _, ok := usedNonceIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for used nonces")
		}
		usedNonceIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in remote token messengers
	tokenMessengerIndexMap := make(map[string]struct{})
	for _, elem := range gs.TokenMessengerList {
		index := string(RemoteTokenMessengerKey(elem.DomainId))
		if _, ok := tokenMessengerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for remote token messengers")
		}
		tokenMessengerIndexMap[index] = struct{}{}
	}

	return nil
}
