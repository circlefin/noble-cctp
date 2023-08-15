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
	"bytes"
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) DepositForBurnWithCaller(goCtx context.Context, msg *types.MsgDepositForBurnWithCaller) (*types.MsgDepositForBurnWithCallerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Destination caller must be nonzero. To allow any destination caller, use DepositForBurn().
	emptyByteArr := make([]byte, types.DestinationCallerLen)
	if len(msg.DestinationCaller) == 0 || bytes.Equal(msg.DestinationCaller, emptyByteArr) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDestinationCaller, "invalid destination caller")
	}

	nonce, err := k.depositForBurn(
		ctx,
		msg.From,
		msg.Amount,
		msg.DestinationDomain,
		msg.MintRecipient,
		msg.BurnToken,
		msg.DestinationCaller)

	return &types.MsgDepositForBurnWithCallerResponse{Nonce: nonce}, err
}
