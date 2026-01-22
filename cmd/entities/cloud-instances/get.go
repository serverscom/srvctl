package cloudinstances

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <instance-id>",
		Short: "Get cloud instance by instance id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			item, err := scClient.CloudComputingInstances.Get(ctx, instanceID)
			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(item)
		},
	}

	return cmd
}
