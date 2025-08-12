package uplinkbandwidths

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64
	var serverModelID int64
	var uplinkModelID int64

	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.BandwidthOption] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Locations.BandwidthOptions(locationID, serverModelID, uplinkModelID)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.BandwidthOption]{},
		&base.BandwidthTypeOption[serverscom.BandwidthOption]{},
	)

	cmd := base.NewListCmd("list", "uplink-bandwidths", factory, cmdContext, opts...)

	cmd.Flags().Int64Var(&locationID, "location-id", 0, "Location ID (required)")
	cmd.Flags().Int64Var(&serverModelID, "server-model-id", 0, "Server Model ID (required)")
	cmd.Flags().Int64Var(&uplinkModelID, "uplink-model-id", 0, "Uplink Model ID (required)")
	_ = cmd.MarkFlagRequired("location-id")
	_ = cmd.MarkFlagRequired("server-model-id")
	_ = cmd.MarkFlagRequired("uplink-model-id")

	return cmd
}
