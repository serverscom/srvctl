package cloudbackups

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddFlags struct {
	VolumeID    string
	Name        string
	Incremental bool
	Force       bool
	Labels      []string
}

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddFlags{}

	cmd := &cobra.Command{
		Use:   "add --volume-id <volume-id> --name <name>",
		Short: "Add a cloud backup",
		Long:  "Create a new cloud block storage backup",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.CloudBlockStorageBackupCreateInput{}

			if err := flags.FillInput(cmd, input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			backup, err := scClient.CloudBlockStorageBackups.Create(ctx, *input)
			if err != nil {
				return err
			}

			if backup != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(backup)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.VolumeID, "volume-id", "", "", "ID of the volume to backup")
	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "Name of the backup")
	cmd.Flags().BoolVarP(&flags.Incremental, "incremental", "", false, "Create incremental backup")
	cmd.Flags().BoolVarP(&flags.Force, "force", "", false, "Force backup creation")
	cmd.Flags().StringArrayVarP(&flags.Labels, "label", "l", []string{}, "string in key=value format")

	_ = cmd.MarkFlagRequired("volume-id")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func (f *AddFlags) FillInput(cmd *cobra.Command, input *serverscom.CloudBlockStorageBackupCreateInput) error {
	input.VolumeID = f.VolumeID
	input.Name = f.Name
	if cmd.Flags().Changed("incremental") {
		input.Incremental = f.Incremental
	}
	if cmd.Flags().Changed("force") {
		input.Force = f.Force
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
