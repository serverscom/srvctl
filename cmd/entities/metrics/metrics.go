package metrics

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "metrics",
		Short:             "Get metrics in Prometheus text exposition format",
		PersistentPreRunE: base.CheckEmptyContexts(cmdContext),
		Args:              base.NoArgs,
		Run:               base.UsageRun,
	}

	cmd.AddCommand(
		newHostsCmd(cmdContext),
		newRacksCmd(cmdContext),
	)

	return cmd
}
