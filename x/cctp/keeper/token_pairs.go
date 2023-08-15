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

// GetTokenPair returns a token pair
func (k Keeper) GetTokenPairHex(ctx sdk.Context, remoteDomain uint32, remoteTokenHex string) (val types.TokenPair, found bool) {
	remoteTokenPadded, err := types.RemoteTokenPadded(remoteTokenHex)
	if err != nil {
		return val, false
	}

	return k.GetTokenPair(ctx, remoteDomain, remoteTokenPadded)
}

// GetTokenPair returns a token pair
func (k Keeper) GetTokenPair(ctx sdk.Context, remoteDomain uint32, remoteToken []byte) (val types.TokenPair, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))

	b := store.Get(types.TokenPairKey(remoteDomain, remoteToken))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetTokenPair sets a token pair in the store
func (k Keeper) SetTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))
	b := k.cdc.MustMarshal(&tokenPair)
	store.Set(types.TokenPairKey(tokenPair.RemoteDomain, tokenPair.RemoteToken), b)
}

// DeleteTokenPair removes a token pair
func (k Keeper) DeleteTokenPair(
	ctx sdk.Context,
	remoteDomain uint32,
	remoteToken []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))
	store.Delete(types.TokenPairKey(remoteDomain, remoteToken))
}

// GetAllTokenPairs returns all token pairs
func (k Keeper) GetAllTokenPairs(ctx sdk.Context) (list []types.TokenPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TokenPair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
