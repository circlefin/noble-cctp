package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgReplaceMessage = "replace_message"

var _ sdk.Msg = &MsgReplaceMessage{}

func NewMsgReplaceMessage(from string, originalMessage []byte, originalAttestation []byte, newMessageBody []byte, newDestinationCaller []byte) *MsgReplaceMessage {
	return &MsgReplaceMessage{
		From:                 from,
		OriginalMessage:      originalMessage,
		OriginalAttestation:  originalAttestation,
		NewMessageBody:       newMessageBody,
		NewDestinationCaller: newDestinationCaller,
	}
}

func (msg *MsgReplaceMessage) Route() string {
	return RouterKey
}

func (msg *MsgReplaceMessage) Type() string {
	return TypeMsgReplaceMessage
}

func (msg *MsgReplaceMessage) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgReplaceMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgReplaceMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
