package sbmosoptions

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64
	var sbmFlavorModelID int64

	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.OperatingSystemOption] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Locations.SBMOperatingSystemOptions(locationID, sbmFlavorModelID)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.OperatingSystemOption]{},
	)

	cmd := base.NewListCmd("list", "sbm-os-options", factory, cmdContext, opts...)

	cmd.Flags().Int64Var(&locationID, "location-id", 0, "Location ID (required)")
	cmd.Flags().Int64Var(&sbmFlavorModelID, "model-id", 0, "SBM flavor model id (int, required)")
	_ = cmd.MarkFlagRequired("location-id")
	_ = cmd.MarkFlagRequired("model-id")

	return cmd
}
