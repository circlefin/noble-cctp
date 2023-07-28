package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTokenMessenger returns a token messenger
func (k Keeper) GetTokenMessenger(ctx sdk.Context, remoteDomain uint32) (val types.TokenMessenger, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenMessengerKeyPrefix))

	b := store.Get(types.TokenMessengerKey(remoteDomain))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetTokenMessenger sets a token messenger in the store
func (k Keeper) SetTokenMessenger(ctx sdk.Context, tokenMessenger types.TokenMessenger) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenMessengerKeyPrefix))
	b := k.cdc.MustMarshal(&tokenMessenger)
	store.Set(types.TokenMessengerKey(tokenMessenger.DomainId), b)
}

// DeleteTokenMessenger removes a token messenger
func (k Keeper) DeleteTokenMessenger(
	ctx sdk.Context,
	remoteDomain uint32,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenMessengerKeyPrefix))
	store.Delete(types.TokenMessengerKey(remoteDomain))
}

// GetAllTokenMessengers returns all token messengers
func (k Keeper) GetAllTokenMessengers(ctx sdk.Context) (list []types.TokenMessenger) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenMessengerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TokenMessenger
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
