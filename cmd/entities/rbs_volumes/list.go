package rbsvolumes

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.RemoteBlockStorageVolume] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.RemoteBlockStorageVolumes.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.RemoteBlockStorageVolume]{},
		&base.LabelSelectorOption[serverscom.RemoteBlockStorageVolume]{},
	)

	return base.NewListCmd("list", "RBS servers", factory, cmdContext, opts...)
}
