package cli

import (
	"strconv"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdSendMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-message [destination-domain] [recipient] [message-body]",
		Short: "Broadcast message send-message",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			destinationDomain, err := strconv.ParseUint(args[0], types.BaseTen, types.DomainBitLen)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendMessage(
				clientCtx.GetFromAddress().String(),
				uint32(destinationDomain),
				[]byte(args[1]),
				[]byte(args[2]),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
