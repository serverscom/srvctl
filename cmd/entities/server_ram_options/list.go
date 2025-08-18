package serverramoptions

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64
	var serverModelID int64

	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.RAMOption] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Locations.RAMOptions(locationID, serverModelID)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.RAMOption]{},
	)

	cmd := base.NewListCmd("list", "server-ram-options", factory, cmdContext, opts...)

	cmd.Flags().Int64Var(&locationID, "location-id", 0, "Location ID (required)")
	cmd.Flags().Int64Var(&serverModelID, "server-model-id", 0, "Server model ID (required)")
	_ = cmd.MarkFlagRequired("location-id")
	_ = cmd.MarkFlagRequired("server-model-id")

	return cmd
}
