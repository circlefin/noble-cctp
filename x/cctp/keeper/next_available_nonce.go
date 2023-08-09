package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetNextAvailableNonce returns the next available nonce
func (k Keeper) GetNextAvailableNonce(ctx sdk.Context) (val types.Nonce, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NextAvailableNonceKey))

	b := store.Get(types.KeyPrefix(types.NextAvailableNonceKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetNextAvailableNonce sets the next available nonce in the store
func (k Keeper) SetNextAvailableNonce(ctx sdk.Context, key types.Nonce) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NextAvailableNonceKey))
	b := k.cdc.MustMarshal(&key)
	store.Set(types.KeyPrefix(types.NextAvailableNonceKey), b)
}

func (k Keeper) ReserveAndIncrementNonce(ctx sdk.Context) (val types.Nonce) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NextAvailableNonceKey))
	b := store.Get(types.KeyPrefix(types.NextAvailableNonceKey))
	k.cdc.MustUnmarshal(b, &val)

	newNonce := types.Nonce{Nonce: val.Nonce + 1}
	b = k.cdc.MustMarshal(&newNonce)

	store.Set(types.KeyPrefix(types.NextAvailableNonceKey), b)
	return val
}
