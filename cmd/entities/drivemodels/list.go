package drivemodels

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	var locationID int64
	var serverModelID int64

	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DriveModel] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Locations.DriveModelOptions(locationID, serverModelID)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.DriveModel]{},
		&base.SearchPatternOption[serverscom.DriveModel]{},
		&base.DriveMediaTypeOption[serverscom.DriveModel]{},
		&base.DriveInterfaceOption[serverscom.DriveModel]{},
	)

	cmd := base.NewListCmd("list", "drive-models", factory, cmdContext, opts...)

	cmd.Flags().Int64Var(&locationID, "location-id", 0, "Location ID (required)")
	cmd.Flags().Int64Var(&serverModelID, "server-model-id", 0, "Server model ID (required)")
	_ = cmd.MarkFlagRequired("location-id")
	_ = cmd.MarkFlagRequired("server-model-id")

	return cmd
}
