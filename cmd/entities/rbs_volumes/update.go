package rbsvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var name string
	var size int64
	var labels []string

	cmd := &cobra.Command{
		Use:   "update <volume-id>",
		Short: "Update a RBS volume",
		Long:  "Update a Remote Block Storage volume by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap, err := base.ParseLabels(labels)
			if err != nil {
				return err
			}

			input := serverscom.RemoteBlockStorageVolumeUpdateInput{
				Name:   name,
				Size:   size,
				Labels: labelsMap,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			volumeID := args[0]
			volume, err := scClient.RemoteBlockStorageVolumes.Update(ctx, volumeID, input)
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

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the RBS volume")
	cmd.Flags().Int64Var(&size, "size", 0, "size of the volume in GB")
	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}
