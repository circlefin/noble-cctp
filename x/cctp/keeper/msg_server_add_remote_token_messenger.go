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

func (k msgServer) AddRemoteTokenMessenger(goCtx context.Context, msg *types.MsgAddRemoteTokenMessenger) (*types.MsgAddRemoteTokenMessengerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot add remote token messengers")
	}

	_, found := k.GetRemoteTokenMessenger(ctx, msg.DomainId)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrRemoteTokenMessengerAlreadyFound, "a remote token messenger for this domain already exists")
	}

	if len(msg.Address) != 32 {
		return nil, sdkerrors.ErrInvalidAddress
	}

	newRemoteTokenMessenger := types.RemoteTokenMessenger{
		DomainId: msg.DomainId,
		Address:  msg.Address,
	}
	k.SetRemoteTokenMessenger(ctx, newRemoteTokenMessenger)

	event := types.RemoteTokenMessengerAdded{
		Domain:               msg.DomainId,
		RemoteTokenMessenger: msg.Address,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgAddRemoteTokenMessengerResponse{}, err
}
