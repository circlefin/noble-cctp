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

// GetTokenPair returns a token pair
func (k Keeper) GetTokenPairHex(ctx context.Context, remoteDomain uint32, remoteTokenHex string) (val types.TokenPair, found bool) {
	remoteTokenPadded, err := types.RemoteTokenPadded(remoteTokenHex)
	if err != nil {
		return val, false
	}

	return k.GetTokenPair(ctx, remoteDomain, remoteTokenPadded)
}

// GetTokenPair returns a token pair
func (k Keeper) GetTokenPair(ctx context.Context, remoteDomain uint32, remoteToken []byte) (val types.TokenPair, found bool) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.TokenPairKeyPrefix))

	b := store.Get(types.TokenPairKey(remoteDomain, remoteToken))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetTokenPair sets a token pair in the store
func (k Keeper) SetTokenPair(ctx context.Context, tokenPair types.TokenPair) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.TokenPairKeyPrefix))
	b := k.cdc.MustMarshal(&tokenPair)
	store.Set(types.TokenPairKey(tokenPair.RemoteDomain, tokenPair.RemoteToken), b)
}

// DeleteTokenPair removes a token pair
func (k Keeper) DeleteTokenPair(
	ctx context.Context,
	remoteDomain uint32,
	remoteToken []byte,
) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.TokenPairKeyPrefix))
	store.Delete(types.TokenPairKey(remoteDomain, remoteToken))
}

// GetAllTokenPairs returns all token pairs
func (k Keeper) GetAllTokenPairs(ctx context.Context) (list []types.TokenPair) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.KeyPrefix(types.TokenPairKeyPrefix))
	iterator := store.Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TokenPair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
