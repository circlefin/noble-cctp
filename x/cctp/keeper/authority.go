package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAuthority returns the authority
func (k Keeper) GetAuthority(ctx sdk.Context) (val types.Authority, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuthorityKey))

	b := store.Get(types.KeyPrefix(types.AuthorityKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetAuthority set authority in the store
func (k Keeper) SetAuthority(ctx sdk.Context, authority types.Authority) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AuthorityKey))

	b := k.cdc.MustMarshal(&authority)
	store.Set(types.KeyPrefix(types.AuthorityKey), b)
}
