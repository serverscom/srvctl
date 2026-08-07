package metrics

import (
	"fmt"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newHostsCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "Get hosts metrics in Prometheus text exposition format",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			metrics, err := scClient.Metrics.ListHostsMetrics(ctx)
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(cmd.OutOrStdout(), metrics)
			return err
		},
	}

	return cmd
}
