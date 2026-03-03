package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListEBMFeaturesCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DedicatedServerFeature] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Hosts.DedicatedServerFeatures(args[0])
	}

	opts := &base.BaseListOptions[serverscom.DedicatedServerFeature]{}

	cmd := base.NewListCmd("list-features", "Enterprise bare metal server features", factory, cmdContext, opts)
	cmd.Use = "list-features <id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
