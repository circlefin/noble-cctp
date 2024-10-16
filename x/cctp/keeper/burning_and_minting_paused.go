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
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// GetBurningAndMintingPaused returns BurningAndMintingPaused
func (k Keeper) GetBurningAndMintingPaused(ctx context.Context) (val types.BurningAndMintingPaused, found bool) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.BurningAndMintingPausedKey))
	b := store.Get(types.KeyPrefix(types.BurningAndMintingPausedKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetBurningAndMintingPaused set BurningAndMintingPaused in the store
func (k Keeper) SetBurningAndMintingPaused(ctx context.Context, paused types.BurningAndMintingPaused) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.BurningAndMintingPausedKey))
	b := k.cdc.MustMarshal(&paused)
	store.Set(types.KeyPrefix(types.BurningAndMintingPausedKey), b)
}
