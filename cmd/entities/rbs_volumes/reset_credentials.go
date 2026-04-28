package rbsvolumes

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newResetCredentialsCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-credentials <volume-id>",
		Short: "Reset credentials for a RBS volume",
		Long:  "Reset credentials for a Remote Block Storage volume by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			volume, err := scClient.RemoteBlockStorageVolumes.ResetCredentials(ctx, volumeID)
			if err != nil {
				return err
			}

			if volume != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(volume)
			}
			return nil
		},
	}

	return cmd
}
