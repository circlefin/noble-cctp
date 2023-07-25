package types

import (
	"testing"

	"github.com/strangelove-ventures/noble/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgReceiveMessage_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgReceiveMessage
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgReceiveMessage{
				From:        "invalid_address",
				Message:     []byte{1, 2, 3},
				Attestation: []byte{1, 2, 3},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgReceiveMessage{
				From:        sample.AccAddress(),
				Message:     []byte{1, 2, 3},
				Attestation: []byte{1, 2, 3},
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
