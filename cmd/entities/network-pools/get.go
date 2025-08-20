package networkpools

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a network pool",
		Long:  "Get a network pool by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			networkPoolID := args[0]
			model, err := scClient.NetworkPools.Get(ctx, networkPoolID)
			if err != nil {
				return err
			}

			if model != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(model)
			}
			return nil
		},
	}

	return cmd
}

func newGetSubnetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var subnetID string

	cmd := &cobra.Command{
		Use:   "get-subnet <id>",
		Short: "Get a subnetwork",
		Long:  "Get a subnetwork for a network pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			networkPoolID := args[0]
			model, err := scClient.NetworkPools.GetSubnetwork(ctx, networkPoolID, subnetID)
			if err != nil {
				return err
			}

			if model != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(model)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&subnetID, "network-id", "", "Subnetwork id (string, required)")
	_ = cmd.MarkFlagRequired("network-id")

	return cmd
}
