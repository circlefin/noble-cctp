package cli

import (
	"strconv"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdUnlinkTokenPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlink-token-pair [local-token] [remote-token] [remote-domain]",
		Short: "Broadcast message unlink-token-pair",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			remoteDomain, err := strconv.ParseUint(args[2], types.BaseTen, types.DomainBitLen)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUnlinkTokenPair(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
				uint32(remoteDomain),
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
