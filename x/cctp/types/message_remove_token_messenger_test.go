package types

import (
	"testing"

	"github.com/strangelove-ventures/noble/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRemoveTokenMessenger_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveTokenMessenger
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgRemoveTokenMessenger{
				From:     "invalid_address",
				DomainId: uint32(1),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgRemoveTokenMessenger{
				From:     sample.AccAddress(),
				DomainId: uint32(123),
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
