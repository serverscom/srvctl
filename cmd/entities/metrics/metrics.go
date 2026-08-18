package metrics

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Get hosts and racks metrics",
		Long: "Get hosts and racks metrics.\n\n" +
			"With the default text output metrics are folded into a table with one row per host or rack.\n" +
			"Use --output raw to get metrics in the Prometheus text exposition format as returned by the API,\n" +
			"e.g. to feed a Prometheus textfile collector.",
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
