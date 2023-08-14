package cli

import (
	"cosmossdk.io/math"
	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

func CmdUpdatePerMessageBurnLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-per-message-burn-limit [denom] [amount]",
		Short: "Broadcast message update-per-message-burn-limit",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			amount, ok := math.NewIntFromString(args[1])
			if !ok {
				return sdkerrors.Wrapf(types.ErrInvalidAmount, "invalid amount")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePerMessageBurnLimit(
				clientCtx.GetFromAddress().String(),
				args[0],
				amount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
