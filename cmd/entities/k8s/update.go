package k8s

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var labels []string

	cmd := &cobra.Command{
		Use:   "update <cluster-id>",
		Short: "Update a kubernetes cluster",
		Long:  "Update a kubernetes cluster by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap, err := base.ParseLabels(labels)
			if err != nil {
				log.Fatal(err)
			}
			input := serverscom.KubernetesClusterUpdateInput{
				Labels: labelsMap,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			k8sCluster, err := scClient.KubernetesClusters.Update(ctx, id, input)
			if err != nil {
				return err
			}

			if k8sCluster != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(k8sCluster)
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}
