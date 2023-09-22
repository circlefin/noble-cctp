/*
 * Copyright (c) 2023, © Circle Internet Financial, LTD.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cli

import (
	"fmt"
	"strconv"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdSendMessageWithCaller() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-message-with-caller [destination-domain] [recipient] [message-body] [destination-caller]",
		Short: "Send a Message With Caller",
		Long:  "Broadcast a transaction that sends a message with a caller to a provided domain.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			destinationDomain, err := strconv.ParseUint(args[0], types.BaseTen, types.DomainBitLen)
			if err != nil {
				return err
			}

			recipient, err := parseAddress(args[1])
			if err != nil {
				return fmt.Errorf("invalid recipient: %w", err)
			}

			destinationCaller, err := parseAddress(args[3])
			if err != nil {
				return fmt.Errorf("invalid destination caller: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendMessageWithCaller(
				clientCtx.GetFromAddress().String(),
				uint32(destinationDomain),
				recipient,
				[]byte(args[2]),
				destinationCaller,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
