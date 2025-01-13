package context

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	var showDefault bool
	var showNonDefault bool

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			contexts := manager.GetContexts()
			defaultContext := manager.GetDefaultContextName()

			if showDefault {
				contexts = output.FilterDefaultContexts(contexts, defaultContext, true)
			} else if showNonDefault {
				contexts = output.FilterDefaultContexts(contexts, defaultContext, false)
			}

			formatter := output.NewFormatter(cmd.OutOrStdout())
			return formatter.FormatContexts(contexts, defaultContext)
		},
	}

	cmd.Flags().BoolVarP(&showDefault, "default", "d", false, "Show only default context")
	cmd.Flags().BoolVar(&showNonDefault, "no-default", false, "Show only non-default contexts")

	return cmd
}
