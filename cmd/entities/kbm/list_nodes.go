package kbm

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.KubernetesBaremetalNode] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.Hosts.ListKubernetesBaremetalNodes()

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.KubernetesBaremetalNode]{},
		&base.LocationIDOption[serverscom.KubernetesBaremetalNode]{},
		&base.RackIDOption[serverscom.KubernetesBaremetalNode]{},
		&base.LabelSelectorOption[serverscom.KubernetesBaremetalNode]{},
		&base.SearchPatternOption[serverscom.KubernetesBaremetalNode]{},
	)

	return base.NewListCmd("list-nodes", "kubernetes baremetal nodes", factory, cmdContext, opts...)
}
