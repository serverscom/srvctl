package kbm

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetPowerFeedsCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-node-feeds <id>",
		Short: "List power feeds for a KBM node",
		Long:  "List power feeds for a KBM node by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]

			powerFeeds, err := scClient.Hosts.KubernetesBaremetalNodePowerFeeds(ctx, id)
			if err != nil {
				return err
			}

			if powerFeeds != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(powerFeeds)
			}
			return nil
		},
	}

	return cmd
}
