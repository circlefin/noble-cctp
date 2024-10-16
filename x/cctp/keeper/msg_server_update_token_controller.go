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

	"github.com/circlefin/noble-cctp/x/cctp/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateTokenController(goCtx context.Context, msg *types.MsgUpdateTokenController) (*types.MsgUpdateTokenControllerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	currentOwner := k.GetOwner(ctx)
	if currentOwner != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "this message sender cannot update the authority")
	}

	_, err := sdk.AccAddressFromBech32(msg.NewTokenController)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidAddress, "invalid new token controller address (%s)", err)
	}

	currentTokenController := k.GetTokenController(ctx)
	k.SetTokenController(ctx, msg.NewTokenController)

	event := types.TokenControllerUpdated{
		PreviousTokenController: currentTokenController,
		NewTokenController:      msg.NewTokenController,
	}

	err = ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUpdateTokenControllerResponse{}, err
}
