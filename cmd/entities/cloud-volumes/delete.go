package cloudvolumes

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newDeleteCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <volume-id>",
		Short: "Delete a cloud volume",
		Long:  "Delete a cloud volume by volume ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			_, err := scClient.CloudBlockStorageVolumes.Delete(ctx, volumeID)

			return err
		},
	}

	return cmd
}
