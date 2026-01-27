package cloudvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddFlags struct {
	InputPath        string
	Name             string
	RegionID         int
	Size             int
	Description      string
	ImageID          string
	SnapshotID       string
	AttachInstanceID string
	BackupID         string
	Labels           []string
}

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddFlags{}

	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: "Add a cloud volume",
		Long:  "Add a new cloud volume to a Cloud region",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.CloudBlockStorageVolumeCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else {
				required := []string{"name", "region-id"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			if err := flags.FillInput(cmd, input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			volume, err := scClient.CloudBlockStorageVolumes.Create(ctx, *input)
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

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "A name of the cloud volume")
	cmd.Flags().IntVar(&flags.RegionID, "region-id", 0, "ID of the cloud region")
	cmd.Flags().IntVar(&flags.Size, "size", 0, "Size of the volume in GB")
	cmd.Flags().StringVar(&flags.Description, "description", "", "Description of the volume")
	cmd.Flags().StringVar(&flags.ImageID, "image-id", "", "ID of the image to create volume from")
	cmd.Flags().StringVar(&flags.SnapshotID, "snapshot-id", "", "ID of the snapshot to create volume from")
	cmd.Flags().StringVar(&flags.AttachInstanceID, "attach-instance-id", "", "ID of the instance to attach volume to")
	cmd.Flags().StringVar(&flags.BackupID, "backup-id", "", "ID of the backup to create volume from")
	cmd.Flags().StringArrayVarP(&flags.Labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func (f *AddFlags) FillInput(cmd *cobra.Command, input *serverscom.CloudBlockStorageVolumeCreateInput) error {
	if cmd.Flags().Changed("name") {
		input.Name = f.Name
	}
	if cmd.Flags().Changed("region-id") {
		input.RegionID = f.RegionID
	}
	if cmd.Flags().Changed("size") {
		input.Size = f.Size
	}
	if cmd.Flags().Changed("description") {
		input.Description = f.Description
	}
	if cmd.Flags().Changed("image-id") {
		input.ImageID = f.ImageID
	}
	if cmd.Flags().Changed("snapshot-id") {
		input.SnapshotID = f.SnapshotID
	}
	if cmd.Flags().Changed("attach-instance-id") {
		input.AttachInstanceID = f.AttachInstanceID
	}
	if cmd.Flags().Changed("backup-id") {
		input.BackupID = f.BackupID
	}
	if cmd.Flags().Changed("label") {
		labelsMap, err := base.ParseLabels(f.Labels)
		if err != nil {
			return err
		}
		input.Labels = labelsMap
	}

	return nil
}
