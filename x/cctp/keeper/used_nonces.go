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
