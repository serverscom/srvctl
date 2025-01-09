package context

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/validator"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var name string
	var setDefault bool

	cmd := &cobra.Command{
		Use:   "update <context-name>",
		Short: "Update a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			contextName := args[0]
			ctx, err := manager.GetContext(contextName)
			if err != nil {
				return err
			}

			if name != "" {
				if err := validator.ValidateContextName(name); err != nil {
					return err
				}
				oldName := ctx.Name
				ctx.Name = name
				if oldName == manager.GetDefaultContextName() {
					manager.SetDefaultContext(name)
				}
			}

			if setDefault {
				if err := manager.SetDefaultContext(ctx.Name); err != nil {
					return err
				}
			}

			return manager.Save()
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "New context name")
	cmd.Flags().BoolVarP(&setDefault, "default", "d", false, "Set as default context")

	return cmd
}
