package cloudregions

import (
	"strconv"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListSnapshotsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.CloudSnapshot] {
		regionID, _ := strconv.ParseInt(args[0], 10, 64)
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.CloudComputingRegions.Snapshots(regionID)

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.CloudSnapshot]{},
		&base.InstanceIDOption[serverscom.CloudSnapshot]{},
		&base.IsBackupOption[serverscom.CloudSnapshot]{},
	)

	cmd := base.NewListCmd("list-snapshots", "cloud region snapshots", factory, cmdContext, opts...)
	cmd.Use = "list-snapshots <region-id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
