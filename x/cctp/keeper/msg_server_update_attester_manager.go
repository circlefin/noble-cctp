/*
 * Copyright (c) 2023, © Circle Internet Financial, LTD.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package keeper

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateAttesterManager(goCtx context.Context, msg *types.MsgUpdateAttesterManager) (*types.MsgUpdateAttesterManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	currentOwner := k.GetOwner(ctx)
	if currentOwner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot update the attester manager")
	}

	currentAttesterManager := k.GetAttesterManager(ctx)
	k.SetAttesterManager(ctx, msg.NewAttesterManager)

	event := types.AttesterManagerUpdated{
		PreviousAttesterManager: currentAttesterManager,
		NewAttesterManager:      msg.NewAttesterManager,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUpdateAttesterManagerResponse{}, err
}
