package config

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

var (
	KnownConfigFlags = []string{"proxy", "http-timeout", "verbose", "output"}
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "config",
		Short:             "Manage configuration",
		Long:              `Manage global and context-specific configurations`,
		PersistentPreRunE: base.CheckEmptyContexts(cmdContext),
	}

	globalCmd := &cobra.Command{
		Use: "global",
	}
	contextCmd := &cobra.Command{
		Use: "context",
	}

	globalCmd.AddCommand(newUpdateCmd(cmdContext))
	contextCmd.AddCommand(newUpdateCmd(cmdContext))

	cmd.AddCommand(
		newFinalCmd(cmdContext),
		globalCmd,
		contextCmd,
	)

	return cmd
}
