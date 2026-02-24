package cloudvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newVolumeAttachCmd(cmdContext *base.CmdContext) *cobra.Command {
	var instanceID string

	cmd := &cobra.Command{
		Use:   "volume-attach <volume-id>",
		Short: "Attach cloud volume to cloud instance",
		Long:  "Attach a cloud volume to a cloud instance",
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

			input := serverscom.CloudBlockStorageVolumeAttachInput{
				InstanceID: instanceID,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			volume, err := scClient.CloudBlockStorageVolumes.Attach(ctx, volumeID, input)
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
