package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnlinkTokenPair = "unlink_token_pair"

var _ sdk.Msg = &MsgUnlinkTokenPair{}

func NewMsgUnlinkTokenPair(from string, localToken string, remoteToken string, remoteDomain uint32) *MsgUnlinkTokenPair {
	return &MsgUnlinkTokenPair{
		From:         from,
		LocalToken:   localToken,
		RemoteDomain: remoteDomain,
		RemoteToken:  remoteToken,
	}
}

func (msg *MsgUnlinkTokenPair) Route() string {
	return RouterKey
}

func (msg *MsgUnlinkTokenPair) Type() string {
	return TypeMsgUnlinkTokenPair
}

func (msg *MsgUnlinkTokenPair) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgUnlinkTokenPair) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnlinkTokenPair) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
