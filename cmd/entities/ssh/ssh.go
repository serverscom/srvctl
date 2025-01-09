package ssh

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "ssh",
		Short:             "Manage ssh keys",
		PersistentPreRunE: base.CheckEmptyContexts(cmdContext),
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newAddCmd(cmdContext),
		NewGetCmd(cmdContext),
		newUpdateCmd(cmdContext),
		newDeleteCmd(cmdContext),
	)

	return cmd
}
