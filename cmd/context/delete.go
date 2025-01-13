package context

import (
	"fmt"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newDeleteCmd(cmdContext *base.CmdContext) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <context-name>",
		Short: "Delete a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			contextName := args[0]
			isDefault, err := manager.IsDefaultContext(contextName)
			if err != nil {
				return err
			}

			if isDefault && !force {
				return fmt.Errorf("cannot delete default context without --force flag")
			}

			if err := manager.DeleteContext(contextName); err != nil {
				return err
			}

			return manager.Save()
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion of default context")

	return cmd
}
