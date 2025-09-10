package kbm

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
	"log"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var labels []string

	cmd := &cobra.Command{
		Use:   "update-node <id>",
		Short: "Update a KBM node",
		Long:  "Update a KBM node by id",
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
			input := serverscom.KubernetesBaremetalNodeUpdateInput{
				Labels: labelsMap,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			id := args[0]

			server, err := scClient.Hosts.UpdateKubernetesBaremetalNode(ctx, id, input)
			if err != nil {
				return err
			}

			if server != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(server)
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}
