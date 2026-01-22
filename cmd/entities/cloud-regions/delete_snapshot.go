package cloudregions

import (
	"fmt"
	"log"
	"strconv"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newDeleteSnapshotCmd(cmdContext *base.CmdContext) *cobra.Command {
	var snapshotID string

	cmd := &cobra.Command{
		Use:   "delete-snapshot <region-id>",
		Short: "Delete snapshot for cloud region by Region ID",
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

			err = scClient.CloudComputingRegions.DeleteSnapshot(ctx, regionID, snapshotID)
			if err != nil {
				return err
			}

			fmt.Printf("Snapshot %s deleted successfully\n", snapshotID)
			return nil
		},
	}

	cmd.Flags().StringVar(&snapshotID, "snapshot-id", "", "Snapshot ID")
	if err := cmd.MarkFlagRequired("snapshot-id"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
