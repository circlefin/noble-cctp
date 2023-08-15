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

// GetSendingAndReceivingMessagesPaused returns SendingAndReceivingMessagesPaused
func (k Keeper) GetSendingAndReceivingMessagesPaused(ctx sdk.Context) (val types.SendingAndReceivingMessagesPaused, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey))

	b := store.Get(types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetSendingAndReceivingMessagesPaused sets SendingAndReceivingMessagesPaused in the store
func (k Keeper) SetSendingAndReceivingMessagesPaused(ctx sdk.Context, paused types.SendingAndReceivingMessagesPaused) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey))
	b := k.cdc.MustMarshal(&paused)
	store.Set(types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey), b)
}
