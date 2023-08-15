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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DeletePendingOwner deletes the pending owner of the CCTP module from state.
func (k Keeper) DeletePendingOwner(ctx sdk.Context) {
	ctx.KVStore(k.storeKey).Delete(types.PendingOwnerKey)
}

// GetOwner returns the owner of the CCTP module from state.
func (k Keeper) GetOwner(ctx sdk.Context) (owner string) {
	bz := ctx.KVStore(k.storeKey).Get(types.OwnerKey)
	if bz == nil {
		panic("cctp owner not found in state")
	}

	return string(bz)
}

// GetPendingOwner returns the pending owner of the CCTP module from state.
func (k Keeper) GetPendingOwner(ctx sdk.Context) (pendingOwner string, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.PendingOwnerKey)
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

// GetAttesterManager returns the attester manager of the CCTP module from state.
func (k Keeper) GetAttesterManager(ctx sdk.Context) (attesterManager string) {
	bz := ctx.KVStore(k.storeKey).Get(types.AttesterManagerKey)
	if bz == nil {
		panic("cctp attester manager not found in state")
	}

	return string(bz)
}

// GetPauser returns the pauser of the CCTP module from state.
func (k Keeper) GetPauser(ctx sdk.Context) (pauser string) {
	bz := ctx.KVStore(k.storeKey).Get(types.PauserKey)
	if bz == nil {
		panic("cctp pauser not found in state")
	}

	return string(bz)
}

// GetTokenController returns the token controller of the CCTP module from state.
func (k Keeper) GetTokenController(ctx sdk.Context) (tokenController string) {
	bz := ctx.KVStore(k.storeKey).Get(types.TokenControllerKey)
	if bz == nil {
		panic("cctp token controller not found in state")
	}

	return string(bz)
}

// SetOwner stores the owner of the CCTP module in state.
func (k Keeper) SetOwner(ctx sdk.Context, owner string) {
	bz := []byte(owner)
	ctx.KVStore(k.storeKey).Set(types.OwnerKey, bz)
}

// SetPendingOwner stores the pending owner of the CCTP module in state.
func (k Keeper) SetPendingOwner(ctx sdk.Context, pendingOwner string) {
	bz := []byte(pendingOwner)
	ctx.KVStore(k.storeKey).Set(types.PendingOwnerKey, bz)
}

// SetAttesterManager stores the attester manager of the CCTP module in state.
func (k Keeper) SetAttesterManager(ctx sdk.Context, attesterManager string) {
	bz := []byte(attesterManager)
	ctx.KVStore(k.storeKey).Set(types.AttesterManagerKey, bz)
}

// SetPauser stores the pauser of the CCTP module in state.
func (k Keeper) SetPauser(ctx sdk.Context, pauser string) {
	bz := []byte(pauser)
	ctx.KVStore(k.storeKey).Set(types.PauserKey, bz)
}

// SetTokenController stores the token controller of the CCTP module in state.
func (k Keeper) SetTokenController(ctx sdk.Context, tokenController string) {
	bz := []byte(tokenController)
	ctx.KVStore(k.storeKey).Set(types.TokenControllerKey, bz)
}
