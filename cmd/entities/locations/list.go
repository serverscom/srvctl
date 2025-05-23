package locations

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Location] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Locations.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Location]{},
		&base.SearchPatternOption[serverscom.Location]{},
	)

	return base.NewListCmd("list", "locations", factory, cmdContext, opts...)
}
