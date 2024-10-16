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

package keeper

import (
	"fmt"

	"github.com/circlefin/noble-cctp/x/cctp/types"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		logger       log.Logger
		storeService store.KVStoreService

		bank             types.BankKeeper
		fiattokenfactory types.FiatTokenfactoryKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	logger log.Logger,
	storeService store.KVStoreService,
	bank types.BankKeeper,
	fiattokenfactory types.FiatTokenfactoryKeeper,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		logger:           logger,
		storeService:     storeService,
		bank:             bank,
		fiattokenfactory: fiattokenfactory,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
