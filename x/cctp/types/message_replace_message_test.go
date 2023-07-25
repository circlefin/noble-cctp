package types

import (
	"testing"

	"github.com/strangelove-ventures/noble/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgReplaceMessage_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgReplaceMessage
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgReplaceMessage{
				From:                 "invalid_address",
				OriginalMessage:      []byte{1, 2, 3},
				OriginalAttestation:  []byte{1, 2, 3},
				NewMessageBody:       []byte{1, 2, 3},
				NewDestinationCaller: []byte{1, 2, 3},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgReplaceMessage{
				From:                 sample.AccAddress(),
				OriginalMessage:      []byte{1, 2, 3},
				OriginalAttestation:  []byte{1, 2, 3},
				NewMessageBody:       []byte{1, 2, 3},
				NewDestinationCaller: []byte{1, 2, 3},
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
