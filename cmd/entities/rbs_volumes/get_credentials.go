package rbsvolumes

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCredentialsCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-credentials <volume-id>",
		Short: "Get credentials for a RBS volume",
		Long:  "Get credentials for a Remote Block Storage volume by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			credentials, err := scClient.RemoteBlockStorageVolumes.GetCredentials(ctx, volumeID)
			if err != nil {
				return err
			}

			if credentials != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(credentials)
			}
			return nil
		},
	}

	return cmd
}
