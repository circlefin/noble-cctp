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

// GetPerMessageBurnLimit returns a PerMessageBurnLimit
func (k Keeper) GetPerMessageBurnLimit(ctx sdk.Context, denom string) (val types.PerMessageBurnLimit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PerMessageBurnLimitKeyPrefix))

	b := store.Get(types.KeyPrefix(string(types.PerMessageBurnLimitKey(denom))))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetPerMessageBurnLimit sets a PerMessageBurnLimit in the store
func (k Keeper) SetPerMessageBurnLimit(ctx sdk.Context, limit types.PerMessageBurnLimit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PerMessageBurnLimitKeyPrefix))
	b := k.cdc.MustMarshal(&limit)
	store.Set(types.KeyPrefix(string(types.PerMessageBurnLimitKey(limit.Denom))), b)
}

// GetAllMessageBurnLimit gets all PerMessageBurnLimits from the store
func (k Keeper) GetAllPerMessageBurnLimits(ctx sdk.Context) (list []types.PerMessageBurnLimit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PerMessageBurnLimitKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerMessageBurnLimit
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
