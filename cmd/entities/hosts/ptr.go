package hosts

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListDSPTRCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.PTRRecord] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		if len(args) == 0 {
			log.Fatal("Missing dedicated server ID")
		}
		id := args[0]
		return scClient.Hosts.DedicatedServerPTRRecords(id)
	}

	opts := &base.BaseListOptions[serverscom.PTRRecord]{}

	return base.NewListCmd("list-ptr <id>", "Dedicated server PTR records", factory, cmdContext, opts)
}

// TODO add other PTR methods
