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
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgPauseSendingAndReceivingMessages = "pause_sending_and_receiving_messages"

var _ sdk.Msg = &MsgPauseSendingAndReceivingMessages{}

func NewMsgPauseSendingAndReceivingMessages(from string) *MsgPauseSendingAndReceivingMessages {
	return &MsgPauseSendingAndReceivingMessages{
		From: from,
	}
}

func (msg *MsgPauseSendingAndReceivingMessages) Route() string {
	return RouterKey
}

func (msg *MsgPauseSendingAndReceivingMessages) Type() string {
	return TypeMsgPauseSendingAndReceivingMessages
}

func (msg *MsgPauseSendingAndReceivingMessages) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgPauseSendingAndReceivingMessages) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPauseSendingAndReceivingMessages) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
