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

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group cctp queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListAttesters())
	cmd.AddCommand(CmdListPerMessageBurnLimits())
	cmd.AddCommand(CmdRemoteTokenMessengers())
	cmd.AddCommand(CmdListTokenPairs())
	cmd.AddCommand(CmdListUsedNonces())
	cmd.AddCommand(CmdShowAttester())
	cmd.AddCommand(CmdRoles())
	cmd.AddCommand(CmdShowBurningAndMintingPaused())
	cmd.AddCommand(CmdShowMaxMessageBodySize())
	cmd.AddCommand(CmdShowNextAvailableNonce())
	cmd.AddCommand(CmdShowPerMessageBurnLimit())
	cmd.AddCommand(CmdShowSendingAndReceivingMessagesPaused())
	cmd.AddCommand(CmdShowSignatureThreshold())
	cmd.AddCommand(CmdRemoteTokenMessenger())
	cmd.AddCommand(CmdShowTokenPair())
	cmd.AddCommand(CmdShowUsedNonce())
	cmd.AddCommand(CmdBurnMessageVersion())
	cmd.AddCommand(CmdLocalMessageVersion())
	cmd.AddCommand(CmdLocalDomain())

	return cmd
}
