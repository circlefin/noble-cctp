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

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// DeletePendingOwner deletes the pending owner of the CCTP module from state.
func (k Keeper) DeletePendingOwner(ctx context.Context) {
	runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Delete(types.PendingOwnerKey)
}

// GetOwner returns the owner of the CCTP module from state.
func (k Keeper) GetOwner(ctx context.Context) (owner string) {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.OwnerKey)
	if bz == nil {
		panic("cctp owner not found in state")
	}

	return string(bz)
}

// GetPendingOwner returns the pending owner of the CCTP module from state.
func (k Keeper) GetPendingOwner(ctx context.Context) (pendingOwner string, found bool) {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.PendingOwnerKey)
	if bz == nil {
		return "", false
	}

	return string(bz), true
}

// GetAttesterManager returns the attester manager of the CCTP module from state.
func (k Keeper) GetAttesterManager(ctx context.Context) (attesterManager string) {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.AttesterManagerKey)
	if bz == nil {
		panic("cctp attester manager not found in state")
	}

	return string(bz)
}

// GetPauser returns the pauser of the CCTP module from state.
func (k Keeper) GetPauser(ctx context.Context) (pauser string) {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.PauserKey)
	if bz == nil {
		panic("cctp pauser not found in state")
	}

	return string(bz)
}

// GetTokenController returns the token controller of the CCTP module from state.
func (k Keeper) GetTokenController(ctx context.Context) (tokenController string) {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.TokenControllerKey)
	if bz == nil {
		panic("cctp token controller not found in state")
	}

	return string(bz)
}

// SetOwner stores the owner of the CCTP module in state.
func (k Keeper) SetOwner(ctx context.Context, owner string) {
	bz := []byte(owner)
	runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Set(types.OwnerKey, bz)
}

// SetPendingOwner stores the pending owner of the CCTP module in state.
func (k Keeper) SetPendingOwner(ctx context.Context, pendingOwner string) {
	bz := []byte(pendingOwner)
	runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Set(types.PendingOwnerKey, bz)
}

// SetAttesterManager stores the attester manager of the CCTP module in state.
func (k Keeper) SetAttesterManager(ctx context.Context, attesterManager string) {
	bz := []byte(attesterManager)
	runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Set(types.AttesterManagerKey, bz)
}

// SetPauser stores the pauser of the CCTP module in state.
func (k Keeper) SetPauser(ctx context.Context, pauser string) {
	bz := []byte(pauser)
	runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Set(types.PauserKey, bz)
}

// SetTokenController stores the token controller of the CCTP module in state.
func (k Keeper) SetTokenController(ctx context.Context, tokenController string) {
	bz := []byte(tokenController)
	runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Set(types.TokenControllerKey, bz)
}
