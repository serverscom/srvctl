package sbmmodels

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64

	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.SBMFlavor] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Locations.SBMFlavorOptions(locationID)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.SBMFlavor]{},
		&base.SearchPatternOption[serverscom.SBMFlavor]{},
		&base.SBMFlavorsShowAllOption[serverscom.SBMFlavor]{},
	)

	cmd := base.NewListCmd("list", "sbm-models", factory, cmdContext, opts...)

	cmd.Flags().Int64Var(&locationID, "location-id", 0, "Location ID (required)")
	_ = cmd.MarkFlagRequired("location-id")

	return cmd
}
