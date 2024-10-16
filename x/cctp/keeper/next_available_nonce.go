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

// GetNextAvailableNonce returns the next available nonce
func (k Keeper) GetNextAvailableNonce(ctx context.Context) (val types.Nonce, found bool) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.NextAvailableNonceKey))

	b := store.Get(types.KeyPrefix(types.NextAvailableNonceKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetNextAvailableNonce sets the next available nonce in the store
func (k Keeper) SetNextAvailableNonce(ctx context.Context, key types.Nonce) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.NextAvailableNonceKey))
	b := k.cdc.MustMarshal(&key)
	store.Set(types.KeyPrefix(types.NextAvailableNonceKey), b)
}

func (k Keeper) ReserveAndIncrementNonce(ctx context.Context) (val types.Nonce) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.NextAvailableNonceKey))
	b := store.Get(types.KeyPrefix(types.NextAvailableNonceKey))
	k.cdc.MustUnmarshal(b, &val)

	newNonce := types.Nonce{Nonce: val.Nonce + 1}
	b = k.cdc.MustMarshal(&newNonce)

	store.Set(types.KeyPrefix(types.NextAvailableNonceKey), b)
	return val
}
