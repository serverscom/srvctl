package cloudregions

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.CloudComputingRegion] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.CloudComputingRegions.Collection()

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.CloudComputingRegion]{},
	)

	return base.NewListCmd("list", "cloud regions", factory, cmdContext, opts...)
}
