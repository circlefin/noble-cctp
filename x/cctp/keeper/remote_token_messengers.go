package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetRemoteTokenMessenger returns a remote token messenger
func (k Keeper) GetRemoteTokenMessenger(ctx sdk.Context, remoteDomain uint32) (val types.RemoteTokenMessenger, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RemoteTokenMessengerKeyPrefix))

	b := store.Get(types.RemoteTokenMessengerKey(remoteDomain))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetRemoteTokenMessenger sets a remote token messenger in the store
func (k Keeper) SetRemoteTokenMessenger(ctx sdk.Context, remoteTokenMessenger types.RemoteTokenMessenger) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RemoteTokenMessengerKeyPrefix))
	b := k.cdc.MustMarshal(&remoteTokenMessenger)
	store.Set(types.RemoteTokenMessengerKey(remoteTokenMessenger.DomainId), b)
}

// DeleteRemoteTokenMessenger removes a remote token messenger
func (k Keeper) DeleteRemoteTokenMessenger(
	ctx sdk.Context,
	remoteDomain uint32,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RemoteTokenMessengerKeyPrefix))
	store.Delete(types.RemoteTokenMessengerKey(remoteDomain))
}

// GetRemoteTokenMessengers returns all remote token messengers
func (k Keeper) GetRemoteTokenMessengers(ctx sdk.Context) (list []types.RemoteTokenMessenger) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RemoteTokenMessengerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RemoteTokenMessenger
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
