package context

import (
	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var showDefault bool
	var showNonDefault bool

	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := config.NewManager()
			if err != nil {
				return err
			}

			contexts := manager.GetContexts()
			defaultContext := manager.GetDefaultContextName()

			if showDefault {
				contexts = output.FilterDefaultContexts(contexts, defaultContext, true)
			} else if showNonDefault {
				contexts = output.FilterDefaultContexts(contexts, defaultContext, false)
			}

			return output.FormatContexts(contexts, defaultContext)
		},
	}

	cmd.Flags().BoolVarP(&showDefault, "default", "d", false, "Show only default context")
	cmd.Flags().BoolVar(&showNonDefault, "no-default", false, "Show only non-default contexts")

	return cmd
}
