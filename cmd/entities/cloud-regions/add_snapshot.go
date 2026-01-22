package cloudregions

import (
	"log"
	"strconv"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newAddSnapshotCmd(cmdContext *base.CmdContext) *cobra.Command {
	var name string
	var instanceID string

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
				Name:       name,
				InstanceID: instanceID,
			}

			snapshot, err := scClient.CloudComputingRegions.CreateSnapshot(ctx, regionID, input)
			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(snapshot)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Snapshot name")
	cmd.Flags().StringVar(&instanceID, "instance-id", "", "Instance ID")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}
	if err := cmd.MarkFlagRequired("instance-id"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
