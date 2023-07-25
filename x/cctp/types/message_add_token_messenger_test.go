package types

import (
	"github.com/strangelove-ventures/noble/testutil/sample"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestAddTokenMessenger_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddTokenMessenger
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgAddTokenMessenger{
				From:     "invalid_address",
				DomainId: uint32(123),
				Address:  "123",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgAddTokenMessenger{
				From:     sample.AccAddress(),
				DomainId: uint32(123),
				Address:  "123",
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
