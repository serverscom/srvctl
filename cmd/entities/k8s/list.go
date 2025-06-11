package k8s

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.KubernetesCluster] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.KubernetesClusters.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.KubernetesCluster]{},
		&base.SearchPatternOption[serverscom.KubernetesCluster]{},
		&base.LabelSelectorOption[serverscom.KubernetesCluster]{},
		&base.LocationIDOption[serverscom.KubernetesCluster]{},
	)

	return base.NewListCmd("list", "k8s", factory, cmdContext, opts...)
}

func newListNodesCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.KubernetesClusterNode] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.KubernetesClusters.Nodes(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.KubernetesClusterNode]{},
		&base.SearchPatternOption[serverscom.KubernetesClusterNode]{},
		&base.LabelSelectorOption[serverscom.KubernetesClusterNode]{},
	)

	return base.NewListCmd("list-nodes <cluster-id>", "kubernetes cluster nodes by ID", factory, cmdContext, opts...)
}
