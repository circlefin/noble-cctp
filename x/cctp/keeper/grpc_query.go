package keeper

import (
	"github.com/circlefin/noble-cctp/x/cctp/types"
)

var _ types.QueryServer = Keeper{}
