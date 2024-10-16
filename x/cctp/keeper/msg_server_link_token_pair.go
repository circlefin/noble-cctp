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
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

var remoteTokenNumBytes int = 32

func (k msgServer) LinkTokenPair(goCtx context.Context, msg *types.MsgLinkTokenPair) (*types.MsgLinkTokenPairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenController := k.GetTokenController(ctx)
	if tokenController != msg.From {
		return nil, errors.Wrapf(types.ErrUnauthorized, "this message sender cannot link token pairs")
	}

	if len(msg.RemoteToken) != remoteTokenNumBytes {
		return nil, errors.Wrapf(types.ErrInvalidRemoteToken, "must be a byte%d array", remoteTokenNumBytes)
	}

	// check whether there already exists a mapping for this remote domain/token
	_, found := k.GetTokenPair(ctx, msg.RemoteDomain, msg.RemoteToken)
	if found {
		return nil, errors.Wrapf(
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
