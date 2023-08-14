package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetOwner returns the owner of the CCTP module from state.
func (k Keeper) GetOwner(ctx sdk.Context) (owner string) {
	bz := ctx.KVStore(k.storeKey).Get(types.OwnerKey)
	if bz == nil {
		panic("cctp owner not found in state")
	}

	return string(bz)
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
