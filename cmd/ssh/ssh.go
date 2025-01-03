package ssh

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "ssh",
		Short:            "Manage ssh keys",
		PersistentPreRun: base.CheckEmptyContexts,
	}

	cmd.AddCommand(
		newListCmd(),
		newAddCmd(),
		newGetCmd(),
		newUpdateCmd(),
		newDeleteCmd(),
	)

	return cmd
}
