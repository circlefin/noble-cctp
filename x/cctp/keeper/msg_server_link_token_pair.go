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
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) LinkTokenPair(goCtx context.Context, msg *types.MsgLinkTokenPair) (*types.MsgLinkTokenPairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenController := k.GetTokenController(ctx)
	if tokenController != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot link token pairs")
	}

	// check whether there already exists a mapping for this remote domain/token
	_, found := k.GetTokenPair(ctx, msg.RemoteDomain, msg.RemoteToken)
	if found {
		return nil, sdkerrors.Wrapf(
			types.ErrTokenPairAlreadyFound,
			"Local token for this remote domain + remote token mapping already exists in store")
	}

	newTokenPair := types.TokenPair{
		RemoteDomain: msg.RemoteDomain,
		RemoteToken:  msg.RemoteToken,
		LocalToken:   strings.ToLower(msg.LocalToken),
	}

	k.SetTokenPair(ctx, newTokenPair)

	event := types.TokenPairLinked{
		LocalToken:   newTokenPair.LocalToken,
		RemoteDomain: msg.RemoteDomain,
		RemoteToken:  msg.RemoteToken,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)
	return &types.MsgLinkTokenPairResponse{}, err
}
