package kbm

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
	"log"
)

func newListSlotsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.HostDriveSlot] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		if len(args) == 0 {
			log.Fatal("Missing KBM node ID")
		}
		id := args[0]
		return scClient.Hosts.KubernetesBaremetalNodeDriveSlots(id)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.HostDriveSlot]{},
		&base.SearchPatternOption[serverscom.HostDriveSlot]{},
	)

	return base.NewListCmd("list-node-slots <id>", "KBM node drive slots", factory, cmdContext, opts...)
}
