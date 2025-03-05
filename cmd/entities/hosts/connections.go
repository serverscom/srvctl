package hosts

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListDSConnectionsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.HostConnection] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		if len(args) == 0 {
			log.Fatal("Missing dedicated server ID")
		}
		id := args[0]
		return scClient.Hosts.DedicatedServerConnections(id)
	}

	opts := &base.BaseListOptions[serverscom.HostConnection]{}

	return base.NewListCmd("list-connections <id>", "Dedicated server connections", factory, cmdContext, opts)
}
