package k8s

import (
	"log"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a kubernetes cluster",
		Long:  "Get a kubernetes cluster by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			k8s, err := scClient.KubernetesClusters.Get(ctx, id)
			if err != nil {
				return err
			}

			if k8s != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(k8s)
			}
			return nil
		},
	}

	return cmd
}

func newGetNodeCmd(cmdContext *base.CmdContext) *cobra.Command {
	var clusterID string
	cmd := &cobra.Command{
		Use:   "get-node <node-id> --cluster-id <string>",
		Short: "Get a kubernetes cluster node",
		Long:  "Get a kubernetes cluster node by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			nodeID := args[0]
			k8s, err := scClient.KubernetesClusters.GetNode(ctx, clusterID, nodeID)
			if err != nil {
				return err
			}

			if k8s != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(k8s)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&clusterID, "cluster-id", "", "cluster id")
	if err := cmd.MarkFlagRequired("cluster-id"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
