package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBurningAndMintingPaused returns BurningAndMintingPaused
func (k Keeper) GetBurningAndMintingPaused(ctx sdk.Context) (val types.BurningAndMintingPaused, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BurningAndMintingPausedKey))
	b := store.Get(types.KeyPrefix(types.BurningAndMintingPausedKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetBurningAndMintingPaused set BurningAndMintingPaused in the store
func (k Keeper) SetBurningAndMintingPaused(ctx sdk.Context, paused types.BurningAndMintingPaused) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BurningAndMintingPausedKey))
	b := k.cdc.MustMarshal(&paused)
	store.Set(types.KeyPrefix(types.BurningAndMintingPausedKey), b)
}
