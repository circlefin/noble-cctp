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
package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetUsedNonce returns a nonce
func (k Keeper) GetUsedNonce(ctx sdk.Context, nonce types.Nonce) (found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UsedNonceKeyPrefix))
	return store.Get(types.UsedNonceKey(nonce.Nonce, nonce.SourceDomain)) != nil
}

// SetUsedNonce sets a nonce in the store
func (k Keeper) SetUsedNonce(ctx sdk.Context, nonce types.Nonce) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UsedNonceKeyPrefix))
	b := k.cdc.MustMarshal(&nonce)
	store.Set(types.UsedNonceKey(nonce.Nonce, nonce.SourceDomain), b)
}

// GetAllUsedNonces returns all UsedNonces
func (k Keeper) GetAllUsedNonces(ctx sdk.Context) (list []types.Nonce) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.UsedNonceKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Nonce
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
