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
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RemoveRemoteTokenMessenger(goCtx context.Context, msg *types.MsgRemoveRemoteTokenMessenger) (*types.MsgRemoveRemoteTokenMessengerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot remove remote token messengers")
	}

	existingRemoteTokenMessenger, found := k.GetRemoteTokenMessenger(ctx, msg.DomainId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrRemoteTokenMessengerNotFound, "no remote token messenger was found for this domain")
	}

	k.DeleteRemoteTokenMessenger(ctx, msg.DomainId)

	event := types.RemoteTokenMessengerRemoved{
		Domain:               msg.DomainId,
		RemoteTokenMessenger: existingRemoteTokenMessenger.Address,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgRemoveRemoteTokenMessengerResponse{}, err
}
