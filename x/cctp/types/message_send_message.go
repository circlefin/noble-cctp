package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSendMessage = "send_message"

var _ sdk.Msg = &MsgSendMessage{}

func NewMsgSendMessage(from string, destinationDomain uint32, recipient []byte, messageBody []byte) *MsgSendMessage {
	return &MsgSendMessage{
		From:              from,
		DestinationDomain: destinationDomain,
		Recipient:         recipient,
		MessageBody:       messageBody,
	}
}

func (msg *MsgSendMessage) Route() string {
	return RouterKey
}

func (msg *MsgSendMessage) Type() string {
	return TypeMsgSendMessage
}

func (msg *MsgSendMessage) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgSendMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
