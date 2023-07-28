package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MockRouterKeeper struct{}

func (MockRouterKeeper) HandleMessage(sdk.Context, []byte) error {
	return nil
}
