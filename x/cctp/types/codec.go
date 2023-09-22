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
package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAcceptOwner{}, "cctp/AcceptOwner", nil)
	cdc.RegisterConcrete(&MsgAddRemoteTokenMessenger{}, "cctp/AddRemoteTokenMessenger", nil)
	cdc.RegisterConcrete(&MsgDepositForBurn{}, "cctp/DepositForBurn", nil)
	cdc.RegisterConcrete(&MsgDepositForBurnWithCaller{}, "cctp/DepositForBurnWithCaller", nil)
	cdc.RegisterConcrete(&MsgDisableAttester{}, "cctp/DisableAttester", nil)
	cdc.RegisterConcrete(&MsgEnableAttester{}, "cctp/EnableAttester", nil)
	cdc.RegisterConcrete(&MsgLinkTokenPair{}, "cctp/LinkTokenPair", nil)
	cdc.RegisterConcrete(&MsgPauseBurningAndMinting{}, "cctp/PauseBurningAndMinting", nil)
	cdc.RegisterConcrete(&MsgPauseSendingAndReceivingMessages{}, "cctp/PauseSendingAndReceivingMessages", nil)
	cdc.RegisterConcrete(&MsgReceiveMessage{}, "cctp/ReceiveMessage", nil)
	cdc.RegisterConcrete(&MsgRemoveRemoteTokenMessenger{}, "cctp/RemoveRemoteTokenMessenger", nil)
	cdc.RegisterConcrete(&MsgReplaceDepositForBurn{}, "cctp/ReplaceDepositForBurn", nil)
	cdc.RegisterConcrete(&MsgReplaceMessage{}, "cctp/ReplaceMessage", nil)
	cdc.RegisterConcrete(&MsgSendMessage{}, "cctp/SendMessage", nil)
	cdc.RegisterConcrete(&MsgSendMessageWithCaller{}, "cctp/SendMessageWithCaller", nil)
	cdc.RegisterConcrete(&MsgUnlinkTokenPair{}, "cctp/UnlinkTokenPair", nil)
	cdc.RegisterConcrete(&MsgUnpauseBurningAndMinting{}, "cctp/UnpauseBurningAndMinting", nil)
	cdc.RegisterConcrete(&MsgUnpauseSendingAndReceivingMessages{}, "cctp/UnpauseSendingAndReceivingMessages", nil)
	cdc.RegisterConcrete(&MsgUpdateOwner{}, "cctp/UpdateOwner", nil)
	cdc.RegisterConcrete(&MsgUpdateAttesterManager{}, "cctp/UpdateAttesterManager", nil)
	cdc.RegisterConcrete(&MsgUpdateTokenController{}, "cctp/UpdateTokenController", nil)
	cdc.RegisterConcrete(&MsgUpdatePauser{}, "cctp/UpdatePauser", nil)
	cdc.RegisterConcrete(&MsgUpdateMaxMessageBodySize{}, "cctp/UpdateMaxMessageBodySize", nil)
	cdc.RegisterConcrete(&MsgSetMaxBurnAmountPerMessage{}, "cctp/SetMaxBurnAmountPerMessage", nil)
	cdc.RegisterConcrete(&MsgUpdateSignatureThreshold{}, "cctp/UpdateSignatureThreshold", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAcceptOwner{},
		&MsgAddRemoteTokenMessenger{},
		&MsgDepositForBurn{},
		&MsgDepositForBurnWithCaller{},
		&MsgDisableAttester{},
		&MsgEnableAttester{},
		&MsgLinkTokenPair{},
		&MsgPauseBurningAndMinting{},
		&MsgPauseSendingAndReceivingMessages{},
		&MsgReceiveMessage{},
		&MsgRemoveRemoteTokenMessenger{},
		&MsgReplaceDepositForBurn{},
		&MsgReplaceMessage{},
		&MsgSendMessage{},
		&MsgSendMessageWithCaller{},
		&MsgUnlinkTokenPair{},
		&MsgUnpauseBurningAndMinting{},
		&MsgUnpauseSendingAndReceivingMessages{},
		&MsgUpdateOwner{},
		&MsgUpdateAttesterManager{},
		&MsgUpdateTokenController{},
		&MsgUpdatePauser{},
		&MsgUpdateMaxMessageBodySize{},
		&MsgSetMaxBurnAmountPerMessage{},
		&MsgUpdateSignatureThreshold{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
