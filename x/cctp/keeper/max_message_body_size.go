package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetMaxMessageBodySize returns the MaxMessageBodySize
func (k Keeper) GetMaxMessageBodySize(ctx sdk.Context) (val types.MaxMessageBodySize, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MaxMessageBodySizeKey))

	b := store.Get(types.KeyPrefix(types.MaxMessageBodySizeKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetMaxMessageBodySize sets MaxMessageBodySize in the store
func (k Keeper) SetMaxMessageBodySize(ctx sdk.Context, amount types.MaxMessageBodySize) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MaxMessageBodySizeKey))
	b := k.cdc.MustMarshal(&amount)
	store.Set(types.KeyPrefix(types.MaxMessageBodySizeKey), b)
}
