package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveRemoteTokenMessenger = "remove_remote_token_messenger"

var _ sdk.Msg = &MsgRemoveRemoteTokenMessenger{}

func NewMsgRemoveRemoteTokenMessenger(from string, domainId uint32) *MsgRemoveRemoteTokenMessenger {
	return &MsgRemoveRemoteTokenMessenger{
		From:     from,
		DomainId: domainId,
	}
}

func (msg *MsgRemoveRemoteTokenMessenger) Route() string {
	return RouterKey
}

func (msg *MsgRemoveRemoteTokenMessenger) Type() string {
	return TypeMsgRemoveRemoteTokenMessenger
}

func (msg *MsgRemoveRemoteTokenMessenger) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgRemoveRemoteTokenMessenger) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveRemoteTokenMessenger) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
