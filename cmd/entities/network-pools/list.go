package networkpools

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.NetworkPool] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.NetworkPools.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.NetworkPool]{},
		&base.SearchPatternOption[serverscom.NetworkPool]{},
		&base.LabelSelectorOption[serverscom.NetworkPool]{},
		&base.LocationIDOption[serverscom.NetworkPool]{},
		&base.NetworkPoolTypeOption[serverscom.NetworkPool]{},
	)

	return base.NewListCmd("list", "network pools", factory, cmdContext, opts...)
}

func newListSubnetsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Subnetwork] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.NetworkPools.Subnetworks(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Subnetwork]{},
		&base.SearchPatternOption[serverscom.Subnetwork]{},
		&base.AttachedSubnetworksOption[serverscom.Subnetwork]{},
	)

	return base.NewListCmd("list-subnets", "subnets for a network pool", factory, cmdContext, opts...)
}
