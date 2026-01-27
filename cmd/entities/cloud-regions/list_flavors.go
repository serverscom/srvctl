package cloudregions

import (
	"strconv"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListFlavorsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.CloudComputingFlavor] {
		regionID, _ := strconv.ParseInt(args[0], 10, 64)
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.CloudComputingRegions.Flavors(regionID)

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.CloudComputingFlavor]{},
	)

	cmd := base.NewListCmd("list-flavors", "cloud region flavors", factory, cmdContext, opts...)
	cmd.Use = "list-flavors <region-id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
