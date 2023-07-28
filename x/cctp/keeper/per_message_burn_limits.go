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
