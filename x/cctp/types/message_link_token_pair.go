package types

import (
	"encoding/hex"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgLinkTokenPair = "link_token_pair"

var _ sdk.Msg = &MsgLinkTokenPair{}

func NewMsgLinkTokenPair(from string, localToken string, remoteToken string, remoteDomain uint32) *MsgLinkTokenPair {
	return &MsgLinkTokenPair{
		From:         from,
		LocalToken:   localToken,
		RemoteToken:  remoteToken,
		RemoteDomain: remoteDomain,
	}
}

func (msg *MsgLinkTokenPair) Route() string {
	return RouterKey
}

func (msg *MsgLinkTokenPair) Type() string {
	return TypeMsgLinkTokenPair
}

func (msg *MsgLinkTokenPair) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

func (msg *MsgLinkTokenPair) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLinkTokenPair) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}

	if _, err := hex.DecodeString(strings.TrimPrefix(msg.RemoteToken, "0x")); err != nil {
		return sdkerrors.Wrapf(ErrInvalidRemoteToken, "must be hex string (%s)", err)
	}

	return nil
}
