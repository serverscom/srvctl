package config

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

var (
	KnownConfigFlags = []string{"proxy", "http-timeout", "verbose", "output"}

	globalCmd = &cobra.Command{
		Use: "global",
	}

	contextCmd = &cobra.Command{
		Use: "context",
	}
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "config",
		Short:            "Manage configuration",
		Long:             `Manage global and context-specific configurations`,
		PersistentPreRun: base.CheckEmptyContexts,
	}

	globalCmd.AddCommand(newUpdateCmd())
	contextCmd.AddCommand(newUpdateCmd())

	cmd.AddCommand(
		newFinalCmd(),
		globalCmd,
		contextCmd,
	)

	return cmd
}
