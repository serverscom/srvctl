package cloudinstances

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newRebootCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reboot <instance-id>",
		Short: "Reboot cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			out, err := scClient.CloudComputingInstances.Reboot(ctx, instanceID)
			if err != nil {
				return err
			}

			if out != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(out)
			}
			return nil
		},
	}

	return cmd
}
