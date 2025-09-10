package kbm

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-node <id>",
		Short: "Get a KBM node",
		Long:  "Get a KBM node by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			node, err := scClient.Hosts.GetKubernetesBaremetalNode(ctx, id)
			if err != nil {
				return err
			}

			if node != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(node)
			}
			return nil
		},
	}

	return cmd
}
