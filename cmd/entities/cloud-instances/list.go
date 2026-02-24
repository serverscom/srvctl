package cloudinstances

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.CloudComputingInstance] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.CloudComputingInstances.Collection()

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.CloudComputingInstance]{},
		&base.RegionIDOption[serverscom.CloudComputingInstance]{},
		&base.LocationIDOption[serverscom.CloudComputingInstance]{},
		&base.LabelSelectorOption[serverscom.CloudComputingInstance]{},
	)

	return base.NewListCmd("list", "cloud instances", factory, cmdContext, opts...)
}
