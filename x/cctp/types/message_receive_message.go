package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgReceiveMessage = "receive_message"

var _ sdk.Msg = &MsgReceiveMessage{}

func NewMsgReceiveMessage(from string, message []byte, attestation []byte) *MsgReceiveMessage {
	return &MsgReceiveMessage{
		From:        from,
		Message:     message,
		Attestation: attestation,
	}
}

func (msg *MsgReceiveMessage) Route() string {
	return RouterKey
}

func (msg *MsgReceiveMessage) Type() string {
	return TypeMsgReceiveMessage
}

func (msg *MsgReceiveMessage) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgReceiveMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgReceiveMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
