package servermodels

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64

	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.ServerModelOption] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Locations.ServerModelOptions(locationID)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.ServerModelOption]{},
		&base.SearchPatternOption[serverscom.ServerModelOption]{},
		&base.HasRaidControllerOption[serverscom.ServerModelOption]{},
	)

	cmd := base.NewListCmd("list", "server-models", factory, cmdContext, opts...)

	cmd.Flags().Int64Var(&locationID, "location-id", 0, "Location ID (required)")
	_ = cmd.MarkFlagRequired("location-id")

	return cmd
}
