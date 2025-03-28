package racks

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Rack] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Racks.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Rack]{},
		&base.LabelSelectorOption[serverscom.Rack]{},
		&base.LocationIDOption[serverscom.Rack]{},
	)

	return base.NewListCmd("list", "racks", factory, cmdContext, opts...)
}
