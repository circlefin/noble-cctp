package types

import (
	"testing"

	"github.com/circlefin/noble-cctp/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgLinkTokenPair_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgLinkTokenPair
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgLinkTokenPair{
				From:         "invalid_address",
				RemoteDomain: 1,
				RemoteToken:  "0x12345",
				LocalToken:   "uusdc",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgLinkTokenPair{
				From:         sample.AccAddress(),
				RemoteDomain: 1,
				RemoteToken:  "0x012345",
				LocalToken:   "uusdc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
