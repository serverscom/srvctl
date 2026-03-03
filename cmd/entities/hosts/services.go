package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListEBMServicesCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DedicatedServerService] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.DedicatedServerServices(args[0])
	}

	opts := &base.BaseListOptions[serverscom.DedicatedServerService]{}

	cmd := base.NewListCmd("list-services", "Enterprise bare metal server services", factory, cmdContext, opts)
	cmd.Use = "list-services <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
