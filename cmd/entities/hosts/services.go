package hosts

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListDSServicesCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DedicatedServerService] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		if len(args) == 0 {
			log.Fatal("Missing dedicated server ID")
		}
		id := args[0]
		return scClient.Hosts.DedicatedServerServices(id)
	}

	opts := &base.BaseListOptions[serverscom.DedicatedServerService]{}

	return base.NewListCmd("list-services <id>", "Dedicated server services", factory, cmdContext, opts)
}
