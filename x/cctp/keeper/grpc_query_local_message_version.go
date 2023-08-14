package keeper

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k Keeper) LocalMessageVersion(_ context.Context, _ *types.QueryLocalMessageVersionRequest) (*types.QueryLocalMessageVersionResponse, error) {
	return &types.QueryLocalMessageVersionResponse{Version: types.NobleMessageVersion}, nil
}
