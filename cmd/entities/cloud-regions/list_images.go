package cloudregions

import (
	"strconv"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListImagesCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.CloudComputingImage] {
		regionID, _ := strconv.ParseInt(args[0], 10, 64)
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.CloudComputingRegions.Images(regionID)

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.CloudComputingImage]{},
	)

	cmd := base.NewListCmd("list-images", "cloud region images", factory, cmdContext, opts...)
	cmd.Use = "list-images <region-id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
