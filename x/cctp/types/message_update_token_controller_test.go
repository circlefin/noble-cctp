package types

import (
	"testing"

	"github.com/circlefin/noble-cctp/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateTokenController_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateTokenController
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgUpdateTokenController{
				From:               "invalid_address",
				NewTokenController: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgUpdateTokenController{
				From:               sample.AccAddress(),
				NewTokenController: sample.AccAddress(),
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
