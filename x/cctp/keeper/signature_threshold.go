package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetSignatureThreshold returns the SignatureThreshold
func (k Keeper) GetSignatureThreshold(ctx sdk.Context) (val types.SignatureThreshold, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SignatureThresholdKey))

	b := store.Get(types.KeyPrefix(types.SignatureThresholdKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetSignatureThreshold sets a SignatureThreshold in the store
func (k Keeper) SetSignatureThreshold(ctx sdk.Context, key types.SignatureThreshold) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SignatureThresholdKey))
	b := k.cdc.MustMarshal(&key)
	store.Set(types.KeyPrefix(types.SignatureThresholdKey), b)
}
