package context

import (
	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/validator"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var name string
	var setDefault bool

	cmd := &cobra.Command{
		Use:   "update <context-name>",
		Short: "Update a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			manager, err := config.NewManager()
			if err != nil {
				return err
			}

			ctx, err := manager.GetContext(contextName)
			if err != nil {
				return err
			}

			if name != "" {
				if err := validator.ValidateContextName(name); err != nil {
					return err
				}
				ctx.Name = name
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
