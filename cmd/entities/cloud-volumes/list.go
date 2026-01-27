package cloudvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.CloudBlockStorageVolume] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.CloudBlockStorageVolumes.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.CloudBlockStorageVolume]{},
		&base.LabelSelectorOption[serverscom.CloudBlockStorageVolume]{},
		&base.RegionIDOption[serverscom.CloudBlockStorageVolume]{},
		&base.InstanceIDOption[serverscom.CloudBlockStorageVolume]{},
	)

	return base.NewListCmd("list", "Cloud volumes", factory, cmdContext, opts...)
}
