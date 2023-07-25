package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSendMessageWithCaller = "send_message_with_caller"

var _ sdk.Msg = &MsgSendMessageWithCaller{}

func NewMsgSendMessageWithCaller(from string, destinationDomain uint32, recipient []byte, messageBody []byte, destinationCaller []byte) *MsgSendMessageWithCaller {
	return &MsgSendMessageWithCaller{
		From:              from,
		DestinationDomain: destinationDomain,
		Recipient:         recipient,
		MessageBody:       messageBody,
		DestinationCaller: destinationCaller,
	}
}

func (msg *MsgSendMessageWithCaller) Route() string {
	return RouterKey
}

func (msg *MsgSendMessageWithCaller) Type() string {
	return TypeMsgSendMessageWithCaller
}

func (msg *MsgSendMessageWithCaller) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgSendMessageWithCaller) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendMessageWithCaller) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
