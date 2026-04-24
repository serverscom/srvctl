package rbsvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddFlags struct {
	Skeleton   bool
	InputPath  string
	Name       string
	Size       int64
	LocationID int
	FlavorID   int
	Labels     []string
}

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddFlags{}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a RBS volume",
		Long:  "Add a new Remote Block Storage volume",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("rbs-volumes/add.json")
			}

			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.RemoteBlockStorageVolumeCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else {
				required := []string{"name", "size", "flavor-id"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			if err := flags.FillInput(cmd, input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			volume, err := scClient.RemoteBlockStorageVolumes.Create(ctx, *input)
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
	cmd.Flags().IntVar(&flags.LocationID, "location-id", 0, "ID of the location")
	cmd.Flags().IntVar(&flags.FlavorID, "flavor-id", 0, "ID of the flavor")
	cmd.Flags().StringArrayVarP(&flags.Labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func (f *AddFlags) FillInput(cmd *cobra.Command, input *serverscom.RemoteBlockStorageVolumeCreateInput) error {
	if cmd.Flags().Changed("name") {
		input.Name = f.Name
	}
	if cmd.Flags().Changed("size") {
		input.Size = f.Size
	}
	if cmd.Flags().Changed("location-id") {
		input.LocationID = f.LocationID
	}
	if cmd.Flags().Changed("flavor-id") {
		input.FlavorID = f.FlavorID
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
