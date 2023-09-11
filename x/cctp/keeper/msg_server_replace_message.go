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

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func (k msgServer) ReplaceMessage(goCtx context.Context, msg *types.MsgReplaceMessage) (*types.MsgReplaceMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	paused, found := k.GetSendingAndReceivingMessagesPaused(ctx)
	if found && paused.Paused {
		return nil, sdkerrors.Wrap(types.ErrReplaceMessage, "sending and receiving messages are paused")
	}

	// Validate each signature in the attestation
	// Note: changing attesters or the signature threshold can render all previous messages irreplaceable
	attesters := k.GetAllAttesters(ctx)
	signatureThreshold, found := k.GetSignatureThreshold(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrReplaceMessage, "signature threshold not found")
	}

	if err := VerifyAttestationSignatures(msg.OriginalMessage, msg.OriginalAttestation, attesters, signatureThreshold.Amount); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrSignatureVerification, "unable to verify signatures")
	}

	// validate message format
	originalMessage, err := new(types.Message).Parse(msg.OriginalMessage)
	if err != nil {
		return nil, err
	}

	// validate that the original message sender is the same as this message sender
	messageSender := make([]byte, 32)
	copy(messageSender[12:], sdk.MustAccAddressFromBech32(msg.From))

	if !bytes.Equal(messageSender, originalMessage.Sender) {
		return nil, sdkerrors.Wrap(types.ErrReplaceMessage, "sender not permitted to use nonce")
	}

	// validate source domain
	if originalMessage.SourceDomain != types.NobleDomainId {
		return nil, sdkerrors.Wrap(types.ErrReplaceMessage, "message not originally sent from this domain")
	}

	err = k.sendMessage(
		ctx,
		originalMessage.DestinationDomain,
		originalMessage.Recipient,
		msg.NewDestinationCaller,
		originalMessage.Sender,
		originalMessage.Nonce,
		msg.NewMessageBody)

	return &types.MsgReplaceMessageResponse{}, err
}
