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

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) UpdateMaxMessageBodySize(goCtx context.Context, msg *types.MsgUpdateMaxMessageBodySize) (*types.MsgUpdateMaxMessageBodySizeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := k.GetOwner(ctx)
	if owner != msg.From {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "this message sender cannot update the max message body size")
	}

	newMaxMessageBodySize := types.MaxMessageBodySize{
		Amount: msg.MessageSize,
	}
	k.SetMaxMessageBodySize(ctx, newMaxMessageBodySize)

	event := types.MaxMessageBodySizeUpdated{
		NewMaxMessageBodySize: msg.MessageSize,
	}
	err := ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgUpdateMaxMessageBodySizeResponse{}, err
}
