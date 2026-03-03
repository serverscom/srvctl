package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListDSConnectionsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.HostConnection] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.DedicatedServerConnections(args[0])
	}

	opts := &base.BaseListOptions[serverscom.HostConnection]{}

	cmd := base.NewListCmd("list-connections", "Dedicated server connections", factory, cmdContext, opts)
	cmd.Use = "list-connections <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
