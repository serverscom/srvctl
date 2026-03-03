package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListDSDriveSlotsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.HostDriveSlot] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.DedicatedServerDriveSlots(args[0])
	}

	opts := &base.BaseListOptions[serverscom.HostDriveSlot]{}

	cmd := base.NewListCmd("list-drive-slots", "Dedicated server drive slots", factory, cmdContext, opts)
	cmd.Use = "list-drive-slots <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

func newListKBMDriveSlotsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.HostDriveSlot] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.KubernetesBaremetalNodeDriveSlots(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.HostDriveSlot]{},
		&base.SearchPatternOption[serverscom.HostDriveSlot]{},
	)

	cmd := base.NewListCmd("list-drive-slots", "KBM node drive slots", factory, cmdContext, opts...)
	cmd.Use = "list-drive-slots <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
