package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddRemoteTokenMessenger = "add_remote_token_messenger"

var _ sdk.Msg = &MsgAddRemoteTokenMessenger{}

func NewMsgAddRemoteTokenMessenger(from string, domainId uint32, address string) *MsgAddRemoteTokenMessenger {
	return &MsgAddRemoteTokenMessenger{
		From:     from,
		DomainId: domainId,
		Address:  address,
	}
}

func (msg *MsgAddRemoteTokenMessenger) Route() string {
	return RouterKey
}

func (msg *MsgAddRemoteTokenMessenger) Type() string {
	return TypeMsgAddRemoteTokenMessenger
}

func (msg *MsgAddRemoteTokenMessenger) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgAddRemoteTokenMessenger) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddRemoteTokenMessenger) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
