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
	"encoding/hex"
	"strings"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/crypto"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	fiattokenfactorytypes "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
)

func (k msgServer) DepositForBurn(goCtx context.Context, msg *types.MsgDepositForBurn) (*types.MsgDepositForBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nonce, err := k.depositForBurn(
		ctx,
		msg.From,
		msg.Amount,
		msg.DestinationDomain,
		msg.MintRecipient,
		msg.BurnToken,
		// ([]byte{} here indicates that any address can call receiveMessage()
		// on the destination domain, triggering mint to specified `mintRecipient`)
		[]byte{})

	return &types.MsgDepositForBurnResponse{Nonce: nonce}, err
}

func (k msgServer) depositForBurn(
	ctx sdk.Context,
	from string,
	amount math.Int,
	destinationDomain uint32,
	mintRecipient []byte,
	burnToken string,
	destinationCaller []byte,
) (uint64, error) {
	if !amount.IsPositive() {
		return 0, sdkerrors.Wrap(types.ErrDepositForBurn, "amount must be positive")
	}

	emptyByteArr := make([]byte, types.MintRecipientLen)
	if mintRecipient == nil || bytes.Equal(mintRecipient, emptyByteArr) {
		return 0, sdkerrors.Wrap(types.ErrDepositForBurn, "mint recipient must be nonzero")
	}

	tokenMessenger, found := k.GetRemoteTokenMessenger(ctx, destinationDomain)
	if !found {
		return 0, sdkerrors.Wrap(types.ErrDepositForBurn, "unable to look up destination token messenger")
	}

	// Note: fiat token factory only supports burning 1 token denom
	denom := k.fiattokenfactory.GetMintingDenom(ctx)
	if !strings.EqualFold(denom.Denom, burnToken) {
		return 0, sdkerrors.Wrapf(types.ErrBurn, "burning denom: %s is not supported", burnToken)
	}

	// check if burning/minting is paused
	paused, _ := k.GetBurningAndMintingPaused(ctx)
	if paused.Paused {
		return 0, sdkerrors.Wrap(types.ErrBurn, "burning and minting are paused")
	}

	// check if amount is greater than configured PerMessageBurnLimit for this token
	perMessageBurnLimit, found := k.GetPerMessageBurnLimit(ctx, strings.ToLower(burnToken))
	if found {
		if amount.GT(perMessageBurnLimit.Amount) {
			return 0, sdkerrors.Wrap(types.ErrBurn, "cannot burn more than the maximum per message burn limit")
		}
	}

	// burn coins
	coin := sdk.NewCoin(burnToken, sdk.NewIntFromBigInt(amount.BigInt()))

	err := k.bank.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(from), types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return 0, sdkerrors.Wrap(err, "error during transfer")
	}

	fiatBurnMsg := fiattokenfactorytypes.MsgBurn{
		From:   types.ModuleAddress.String(),
		Amount: coin,
	}
	_, err = k.fiattokenfactory.Burn(ctx, &fiatBurnMsg)
	if err != nil {
		return 0, sdkerrors.Wrapf(err, "error during burn")
	}

	messageSender := make([]byte, 32)
	copy(messageSender[12:], sdk.MustAccAddressFromBech32(from))

	burnMessage := types.BurnMessage{
		Version:       types.MessageBodyVersion,
		BurnToken:     crypto.Keccak256([]byte(strings.ToLower(burnToken))),
		MintRecipient: mintRecipient,
		Amount:        amount,
		MessageSender: messageSender,
	}

	var nonce types.Nonce

	newMessageBodyBytes, err := burnMessage.Bytes()
	if err != nil {
		return 0, sdkerrors.Wrapf(types.ErrParsingBurnMessage, "error parsing burn message into bytes")
	}

	if len(destinationCaller) == 0 {
		message := types.MsgSendMessage{
			From:              types.ModuleAddress.String(),
			DestinationDomain: destinationDomain,
			Recipient:         tokenMessenger.Address,
			MessageBody:       newMessageBodyBytes,
		}

		resp, err := k.SendMessage(sdk.WrapSDKContext(ctx), &message)
		if err != nil {
			return 0, err
		}
		nonce.Nonce = resp.Nonce
	} else {
		message := types.MsgSendMessageWithCaller{
			From:              types.ModuleAddress.String(),
			DestinationDomain: destinationDomain,
			Recipient:         tokenMessenger.Address,
			MessageBody:       newMessageBodyBytes,
			DestinationCaller: destinationCaller,
		}

		resp, err := k.SendMessageWithCaller(sdk.WrapSDKContext(ctx), &message)
		if err != nil {
			return 0, err
		}
		nonce.Nonce = resp.Nonce
	}

	event := types.DepositForBurn{
		Nonce:                     nonce.Nonce,
		BurnToken:                 hex.EncodeToString(crypto.Keccak256([]byte(burnToken))),
		Amount:                    amount,
		Depositor:                 from,
		MintRecipient:             mintRecipient,
		DestinationDomain:         destinationDomain,
		DestinationTokenMessenger: tokenMessenger.Address,
		DestinationCaller:         destinationCaller,
	}
	err = ctx.EventManager().EmitTypedEvent(&event)

	return nonce.Nonce, err
}
