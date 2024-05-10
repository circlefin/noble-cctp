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
	"strings"

	sdkerrors "cosmossdk.io/errors"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	fiattokenfactorytypes "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

var zeroByteArray = []byte{ // 32 bytes
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

func (k msgServer) ReceiveMessage(goCtx context.Context, msg *types.MsgReceiveMessage) (*types.MsgReceiveMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sendingReceivingPaused, found := k.GetSendingAndReceivingMessagesPaused(ctx)
	if found && sendingReceivingPaused.Paused {
		return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "sending and receiving messages are paused")
	}

	// Validate each signature in the attestation
	publicKeys := k.GetAllAttesters(ctx)
	if len(publicKeys) == 0 {
		return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "no attesters found")
	}

	signatureThreshold, found := k.GetSignatureThreshold(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "signature threshold not found")
	}

	if err := VerifyAttestationSignatures(msg.Message, msg.Attestation, publicKeys, signatureThreshold.Amount); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "unable to verify signatures: %s", err)
	}

	// parse message
	message, err := new(types.Message).Parse(msg.Message)
	if err != nil {
		return nil, err
	}

	// validate domain
	if message.DestinationDomain != types.NobleDomainId {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "incorrect destination domain: %d", message.DestinationDomain)
	}

	// validate destination caller
	if !bytes.Equal(message.DestinationCaller, zeroByteArray) {
		bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
		destinationCaller, err := bech32.ConvertAndEncode(bech32Prefix, message.DestinationCaller[12:])
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "unable to encode destination caller %s: %s", msg.From, err)
		}

		if destinationCaller != msg.From {
			return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "incorrect destination caller: %s, sender: %s", destinationCaller, msg.From)
		}
	}

	// validate version
	if message.Version != types.NobleMessageVersion {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "incorrect message version. expected: %d, found: %d", types.NobleMessageVersion, message.Version)
	}

	// validate nonce is available
	// note: we use the domain/nonce combo instead of a hash
	usedNonce := types.Nonce{SourceDomain: message.SourceDomain, Nonce: message.Nonce}
	found = k.GetUsedNonce(ctx, usedNonce)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "nonce already used")
	}

	// mark nonce as used
	k.SetUsedNonce(ctx, usedNonce)

	// verify and parse BurnMessage
	if bytes.Equal(message.Recipient, types.PaddedModuleAddress) { // then mint
		burningMintingPaused, found := k.GetBurningAndMintingPaused(ctx)
		if found && burningMintingPaused.Paused {
			return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "cctp burning and minting is paused")
		}

		burnMessage, err := new(types.BurnMessage).Parse(message.MessageBody)
		if err != nil {
			return nil, err
		}

		if burnMessage.Version != types.MessageBodyVersion {
			return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "invalid message body version")
		}

		// look up Noble mint token from corresponding source domain/token
		tokenPair, found := k.GetTokenPair(ctx, message.SourceDomain, burnMessage.BurnToken)
		if !found {
			return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "corresponding noble mint token not found")
		}

		remoteTokenMessenger, found := k.GetRemoteTokenMessenger(ctx, message.SourceDomain)
		if !found {
			return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "could not retrieve remote token messenger for domain %d", message.SourceDomain)
		}
		if !bytes.Equal(message.Sender, remoteTokenMessenger.Address) {
			return nil, sdkerrors.Wrap(types.ErrReceiveMessage, "message sender is not the remote token messenger")
		}

		// get mint recipient as noble address
		bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
		mintRecipient, err := sdk.Bech32ifyAddressBytes(bech32Prefix, burnMessage.MintRecipient[12:])
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrReceiveMessage, "error bech32 encoding mint recipient address: %s", err)
		}

		msgMint := fiattokenfactorytypes.MsgMint{
			From:    types.ModuleAddress.String(),
			Address: mintRecipient,
			Amount: sdk.Coin{
				Denom:  strings.ToLower(tokenPair.LocalToken),
				Amount: sdk.NewIntFromBigInt(burnMessage.Amount.BigInt()),
			},
		}
		_, err = k.fiattokenfactory.Mint(ctx, &msgMint)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "error during minting: %s", err)
		}

		mintEvent := types.MintAndWithdraw{
			MintRecipient: burnMessage.MintRecipient,
			Amount:        burnMessage.Amount,
			MintToken:     strings.ToLower(tokenPair.LocalToken),
		}
		err = ctx.EventManager().EmitTypedEvent(&mintEvent)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "error emitting mint event: %s", err)
		}
	}

	event := types.MessageReceived{
		Caller:       msg.From,
		SourceDomain: message.SourceDomain,
		Nonce:        message.Nonce,
		Sender:       message.Sender,
		MessageBody:  message.MessageBody,
	}
	err = ctx.EventManager().EmitTypedEvent(&event)

	return &types.MsgReceiveMessageResponse{Success: true}, err
}
