package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetSendingAndReceivingMessagesPaused returns SendingAndReceivingMessagesPaused
func (k Keeper) GetSendingAndReceivingMessagesPaused(ctx sdk.Context) (val types.SendingAndReceivingMessagesPaused, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey))

	b := store.Get(types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetSendingAndReceivingMessagesPaused sets SendingAndReceivingMessagesPaused in the store
func (k Keeper) SetSendingAndReceivingMessagesPaused(ctx sdk.Context, paused types.SendingAndReceivingMessagesPaused) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey))
	b := k.cdc.MustMarshal(&paused)
	store.Set(types.KeyPrefix(types.SendingAndReceivingMessagesPausedKey), b)
}
