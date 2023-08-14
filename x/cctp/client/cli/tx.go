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
	cmd.AddCommand(CmdUpdateAuthority())
	cmd.AddCommand(CmdUpdateMaxMessageBodySize())
	cmd.AddCommand(CmdUpdatePerMessageBurnLimit())
	cmd.AddCommand(CmdUpdateSignatureThreshold())

	return cmd
}
