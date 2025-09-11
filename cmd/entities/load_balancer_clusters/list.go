package loadbalancerclusters

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.LoadBalancerCluster] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.LoadBalancerClusters.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.LoadBalancerCluster]{},
		&base.SearchPatternOption[serverscom.LoadBalancerCluster]{},
	)

	return base.NewListCmd("list", "lb-clusters", factory, cmdContext, opts...)
}
