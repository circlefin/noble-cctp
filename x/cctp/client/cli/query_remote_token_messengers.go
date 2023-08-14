package cli

import (
	"context"
	"strconv"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdRemoteTokenMessengers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote-token-messengers",
		Short: "returns all remote token messengers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRemoteTokenMessengersRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RemoteTokenMessengers(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdRemoteTokenMessenger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote-token-messenger [domain-id]",
		Short: "returns the remote token messenger for a given domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			remoteDomain, _ := strconv.ParseUint(args[0], types.BaseTen, types.DomainBitLen)

			params := &types.QueryRemoteTokenMessengerRequest{
				DomainId: uint32(remoteDomain),
			}

			res, err := queryClient.RemoteTokenMessenger(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
