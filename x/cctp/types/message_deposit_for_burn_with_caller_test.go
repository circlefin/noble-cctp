package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/noble/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgDepositForBurnWithCaller_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDepositForBurnWithCaller
		err  error
	}{
		{
			name: "invalid from",
			msg: MsgDepositForBurnWithCaller{
				From:              "invalid_address",
				Amount:            math.NewInt(123),
				DestinationDomain: 123,
				MintRecipient:     []byte{1, 2, 3},
				BurnToken:         "utoken",
				DestinationCaller: []byte{1, 2, 3},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid from",
			msg: MsgDepositForBurnWithCaller{
				From:              sample.AccAddress(),
				Amount:            math.NewInt(123),
				DestinationDomain: 123,
				MintRecipient:     []byte{1, 2, 3},
				BurnToken:         "utoken",
				DestinationCaller: []byte{1, 2, 3},
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
