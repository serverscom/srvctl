package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Host] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.Hosts.Collection()

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Host]{},
		&base.LocationIDOption[serverscom.Host]{},
		&base.RackIDOption[serverscom.Host]{},
		&base.LabelSelectorOption[serverscom.Host]{},
		&base.SearchPatternOption[serverscom.Host]{},
	)

	return base.NewListCmd("list", "hosts", factory, cmdContext, opts...)
}

func newListDSCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DedicatedServer] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.Hosts.ListDedicatedServers()

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.DedicatedServer]{},
		&base.LocationIDOption[serverscom.DedicatedServer]{},
		&base.RackIDOption[serverscom.DedicatedServer]{},
		&base.LabelSelectorOption[serverscom.DedicatedServer]{},
		&base.SearchPatternOption[serverscom.DedicatedServer]{},
	)

	return base.NewListCmd("list", "dedicated servers", factory, cmdContext, opts...)
}

func newListKBMCmd(cmdContext *base.CmdContext) *cobra.Command {
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

	return base.NewListCmd("list", "dedicated servers", factory, cmdContext, opts...)
}

func newListSBMCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.SBMServer] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		collection := scClient.Hosts.ListSBMServers()

		return collection
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.SBMServer]{},
		&base.LocationIDOption[serverscom.SBMServer]{},
		&base.RackIDOption[serverscom.SBMServer]{},
		&base.LabelSelectorOption[serverscom.SBMServer]{},
		&base.SearchPatternOption[serverscom.SBMServer]{},
	)

	return base.NewListCmd("list", "dedicated servers", factory, cmdContext, opts...)
}
