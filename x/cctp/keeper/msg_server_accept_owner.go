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

func (k msgServer) AcceptOwner(goCtx context.Context, msg *types.MsgAcceptOwner) (*types.MsgAcceptOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	currentOwner := k.GetOwner(ctx)
	pendingOwner, found := k.GetPendingOwner(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "pending owner is not set")
	}

	if pendingOwner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "you are not the pending owner")
	}

	k.SetOwner(ctx, pendingOwner)
	k.DeletePendingOwner(ctx)

	event := types.OwnerUpdated{
		PreviousOwner: currentOwner,
		NewOwner:      pendingOwner,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgAcceptOwnerResponse{}, err
}
