package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAttester returns an attester
func (k Keeper) GetAttester(ctx sdk.Context, key string) (val types.Attester, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttesterKeyPrefix))

	b := store.Get(types.KeyPrefix(string(types.AttesterKey([]byte(key)))))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetAttester sets an attester in the store
func (k Keeper) SetAttester(ctx sdk.Context, key types.Attester) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttesterKeyPrefix))
	b := k.cdc.MustMarshal(&key)
	store.Set(types.KeyPrefix(string(types.AttesterKey([]byte(key.Attester)))), b)
}

// DeleteAttester removes an attester
func (k Keeper) DeleteAttester(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttesterKeyPrefix))
	store.Delete(types.AttesterKey([]byte(key)))
}

// GetAllAttesters returns all attesters
func (k Keeper) GetAllAttesters(ctx sdk.Context) (list []types.Attester) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AttesterKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Attester
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
