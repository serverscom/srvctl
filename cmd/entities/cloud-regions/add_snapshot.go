package cloudregions

import (
	"strconv"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddedFlags struct {
	Name       string
	InstanceID string
}

func newAddSnapshotCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddedFlags{}

	cmd := &cobra.Command{
		Use:   "add-snapshot <region-id>",
		Short: "Add snapshot for cloud region by Region ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			regionID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			input := serverscom.CloudSnapshotCreateInput{
				Name:       flags.Name,
				InstanceID: flags.InstanceID,
			}

			snapshot, err := scClient.CloudComputingRegions.CreateSnapshot(ctx, regionID, input)
			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(snapshot)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "Snapshot name")
	cmd.Flags().StringVar(&flags.InstanceID, "instance-id", "", "Instance ID")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("instance-id")

	return cmd
}
