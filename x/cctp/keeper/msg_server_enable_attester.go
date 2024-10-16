// Copyright 2024 Circle Internet Group, Inc.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) EnableAttester(goCtx context.Context, msg *types.MsgEnableAttester) (*types.MsgEnableAttesterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	attesterManager := k.GetAttesterManager(ctx)
	if attesterManager != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "this message sender cannot enable attesters")
	}

	if len(common.FromHex(msg.Attester)) == 0 {
		return nil, errors.Wrapf(types.ErrInvalidAddress, "invalid attester")
	}

	_, found := k.GetAttester(ctx, msg.Attester)
	if found {
		return nil, errors.Wrapf(types.ErrAttesterAlreadyFound, "this attester already exists in the store")
	}

	newAttester := types.Attester{
		Attester: msg.Attester,
	}
	k.SetAttester(ctx, newAttester)

	event := types.AttesterEnabled{
		Attester: msg.Attester,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgEnableAttesterResponse{}, err
}
