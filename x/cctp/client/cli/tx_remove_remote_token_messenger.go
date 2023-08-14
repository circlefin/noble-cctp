package cli

import (
	"strconv"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdRemoveRemoteTokenMessenger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-remote-token-messenger [domain-id]",
		Short: "Broadcast message remove-remote-token-messenger",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			domainId, err := strconv.ParseUint(args[0], types.BaseTen, types.DomainBitLen)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveRemoteTokenMessenger(
				clientCtx.GetFromAddress().String(),
				uint32(domainId),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
