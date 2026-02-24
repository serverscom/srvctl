package cloudinstances

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpgradeCmd(cmdContext *base.CmdContext) *cobra.Command {
	var flavorID string

	cmd := &cobra.Command{
		Use:   "upgrade <instance-id>",
		Short: "Upgrade cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.CloudComputingInstanceUpgradeInput{
				FlavorID: flavorID,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			out, err := scClient.CloudComputingInstances.Upgrade(ctx, instanceID, input)
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

	cmd.Flags().StringVar(&flavorID, "flavor-id", "", "Flavor ID for upgrade")

	return cmd
}

func newUpgradeApproveCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade-approve <instance-id>",
		Short: "Approve upgrade of cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			out, err := scClient.CloudComputingInstances.ApproveUpgrade(ctx, instanceID)
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

func newUpgradeRevertCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade-revert <instance-id>",
		Short: "Revert upgrade of cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			out, err := scClient.CloudComputingInstances.RevertUpgrade(ctx, instanceID)
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
