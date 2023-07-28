package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTokenPair returns a token pair
func (k Keeper) GetTokenPair(ctx sdk.Context, remoteDomain uint32, remoteToken string) (val types.TokenPair, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))

	b := store.Get(types.TokenPairKey(remoteDomain, remoteToken))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetTokenPair sets a token pair in the store
func (k Keeper) SetTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))
	b := k.cdc.MustMarshal(&tokenPair)
	store.Set(types.TokenPairKey(tokenPair.RemoteDomain, tokenPair.RemoteToken), b)
}

// DeleteTokenPair removes a token pair
func (k Keeper) DeleteTokenPair(
	ctx sdk.Context,
	remoteDomain uint32,
	remoteToken string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))
	store.Delete(types.TokenPairKey(remoteDomain, remoteToken))
}

// GetAllTokenPairs returns all token pairs
func (k Keeper) GetAllTokenPairs(ctx sdk.Context) (list []types.TokenPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenPairKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TokenPair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
