package kbm

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
	"log"
)

func newListNetworksCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Network] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		if len(args) == 0 {
			log.Fatal("Missing KBM node ID")
		}
		id := args[0]
		return scClient.Hosts.KubernetesBaremetalNodeNetworks(id)
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Network]{},
		&base.SearchPatternOption[serverscom.Network]{},
		&base.FamilyOption[serverscom.Network]{},
		&base.InterfaceTypeOption[serverscom.Network]{},
		&base.DistributionMethodOption[serverscom.Network]{},
		&base.AdditionalOption[serverscom.Network]{},
	)

	return base.NewListCmd("list-node-networks <id>", "KBM node networks", factory, cmdContext, opts...)
}
