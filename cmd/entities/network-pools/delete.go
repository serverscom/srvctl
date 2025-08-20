package networkpools

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newDeleteSubnetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var subnetID string

	cmd := &cobra.Command{
		Use:   "delete <network_pool_id>",
		Short: "Delete a subnetwork for a network pool",
		Long:  "Delete a subnetwork for a network pool by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			networkPoolID := args[0]
			return scClient.NetworkPools.DeleteSubnetwork(ctx, networkPoolID, subnetID)
		},
	}

	cmd.Flags().StringVar(&subnetID, "network-id", "", "Subnetwork id (string, required)")
	_ = cmd.MarkFlagRequired("network-id")

	return cmd
}
