package cloudbackups

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.CloudBlockStorageBackup] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.CloudBlockStorageBackups.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.CloudBlockStorageBackup]{},
		&base.LabelSelectorOption[serverscom.CloudBlockStorageBackup]{},
		&base.RegionIDOption[serverscom.CloudBlockStorageBackup]{},
	)

	return base.NewListCmd("list", "Cloud Backups", factory, cmdContext, opts...)
}
