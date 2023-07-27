package cli

import (
	"context"
	"strconv"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListTokenMessengers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-token-messengers",
		Short: "lists all token messengers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllTokenMessengersRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.TokenMessengers(context.Background(), params)
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

func CmdShowTokenMessenger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-token-messenger [domain-id]",
		Short: "shows a token messenger for a given domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			remoteDomain, err := strconv.ParseUint(args[0], types.BaseTen, types.DomainBitLen)

			params := &types.QueryGetTokenMessengerRequest{
				DomainId: uint32(remoteDomain),
			}

			res, err := queryClient.TokenMessenger(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
