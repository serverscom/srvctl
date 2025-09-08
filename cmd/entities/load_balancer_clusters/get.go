package loadbalancerclusters

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get LB Cluster",
		Long:  "Get LB Cluster by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			lbCluster, err := scClient.LoadBalancerClusters.GetLoadBalancerCluster(ctx, id)
			if err != nil {
				return err
			}

			if lbCluster != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(lbCluster)
			}
			return nil
		},
	}

	return cmd
}
