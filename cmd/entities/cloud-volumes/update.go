package cloudvolumes

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var name string
	var description string
	var imageID string
	var snapshotID string
	var labels []string

	cmd := &cobra.Command{
		Use:   "update <volume-id>",
		Short: "Update a cloud volume",
		Long:  "Update a cloud volume by volume ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap, err := base.ParseLabels(labels)
			if err != nil {
				log.Fatal(err)
			}
			input := serverscom.CloudBlockStorageVolumeUpdateInput{
				Name:        name,
				Description: description,
				ImageID:     imageID,
				SnapshotID:  snapshotID,
				Labels:      labelsMap,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			volume, err := scClient.CloudBlockStorageVolumes.Update(ctx, volumeID, input)
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

	cmd.Flags().StringVarP(&name, "name", "n", "", "A name of the cloud volume")
	cmd.Flags().StringVar(&description, "description", "", "Description of the volume")
	cmd.Flags().StringVar(&imageID, "image-id", "", "ID of the image")
	cmd.Flags().StringVar(&snapshotID, "snapshot-id", "", "ID of the snapshot")
	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}
