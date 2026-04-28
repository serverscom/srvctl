package rbsvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type UpdateFlags struct {
	Skeleton  bool
	InputPath string
	Name      string
	Size      int64
	Labels    []string
}

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &UpdateFlags{}

	cmd := &cobra.Command{
		Use:   "update <volume-id>",
		Short: "Update a RBS volume",
		Long:  "Update a Remote Block Storage volume by ID",
		Args:  base.SkeletonOrExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("rbs-volumes/update.json")
			}

			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.RemoteBlockStorageVolumeUpdateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			}

			if err := flags.FillInput(cmd, input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			volume, err := scClient.RemoteBlockStorageVolumes.Update(ctx, volumeID, *input)
			if err != nil {
				return err
			}

			if volume != nil {
				return formatter.Format(volume)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")
	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "name of the RBS volume")
	cmd.Flags().Int64Var(&flags.Size, "size", 0, "size of the volume in GB")
	cmd.Flags().StringArrayVarP(&flags.Labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func (f *UpdateFlags) FillInput(cmd *cobra.Command, input *serverscom.RemoteBlockStorageVolumeUpdateInput) error {
	if cmd.Flags().Changed("name") {
		input.Name = f.Name
	}
	if cmd.Flags().Changed("size") {
		input.Size = f.Size
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
