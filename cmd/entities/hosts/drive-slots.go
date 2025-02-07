package hosts

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListDSDriveSlotsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.HostDriveSlot] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		if len(args) == 0 {
			log.Fatal("Missing dedicated server ID")
		}
		id := args[0]
		return scClient.Hosts.DedicatedServerDriveSlots(id)
	}

	opts := &base.BaseListOptions[serverscom.HostDriveSlot]{}

	return base.NewListCmd("list-drive-slots <id>", "Dedicated server drive slots", factory, cmdContext, opts)
}

// TODO add list drive slots for KBM after adding KBMDriveSlots in go client
