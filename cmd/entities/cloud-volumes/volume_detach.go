package cloudvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newVolumeDetachCmd(cmdContext *base.CmdContext) *cobra.Command {
	var instanceID string

	cmd := &cobra.Command{
		Use:   "volume-detach <volume-id>",
		Short: "Detach cloud volume from cloud instance",
		Long:  "Detach a cloud volume from a cloud instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			required := []string{"instance-id"}
			if err := base.ValidateFlags(cmd, required); err != nil {
				return err
			}

			input := serverscom.CloudBlockStorageVolumeDetachInput{
				InstanceID: instanceID,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			volume, err := scClient.CloudBlockStorageVolumes.Detach(ctx, volumeID, input)
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

	cmd.Flags().StringVar(&instanceID, "instance-id", "", "ID of the cloud instance")

	return cmd
}
