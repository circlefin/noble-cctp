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

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdAcceptOwner())
	cmd.AddCommand(CmdAddRemoteTokenMessenger())
	cmd.AddCommand(CmdDepositForBurn())
	cmd.AddCommand(CmdDepositForBurnWithCaller())
	cmd.AddCommand(CmdDisableAttester())
	cmd.AddCommand(CmdEnableAttester())
	cmd.AddCommand(CmdLinkTokenPair())
	cmd.AddCommand(CmdPauseBurningAndMinting())
	cmd.AddCommand(CmdPauseSendingAndReceivingMessages())
	cmd.AddCommand(CmdReceiveMessage())
	cmd.AddCommand(CmdRemoveRemoteTokenMessenger())
	cmd.AddCommand(CmdReplaceDepositForBurn())
	cmd.AddCommand(CmdReplaceMessage())
	cmd.AddCommand(CmdSendMessage())
	cmd.AddCommand(CmdSendMessageWithCaller())
	cmd.AddCommand(CmdUnlinkTokenPair())
	cmd.AddCommand(CmdUnpauseBurningAndMinting())
	cmd.AddCommand(CmdUnpauseSendingAndReceivingMessages())
	cmd.AddCommand(CmdUpdateOwner())
	cmd.AddCommand(CmdUpdateMaxMessageBodySize())
	cmd.AddCommand(CmdUpdateMaxBurnAmountPerMessage())
	cmd.AddCommand(CmdUpdateSignatureThreshold())
	cmd.AddCommand(CmdUpdateAttesterManager())
	cmd.AddCommand(CmdUpdateTokenController())
	cmd.AddCommand(CmdUpdatePauser())

	return cmd
}
