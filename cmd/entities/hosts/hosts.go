package hosts

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

type HostTypeCmd struct {
	use        string
	shortDesc  string
	entityName string
	typeFlag   string
	managers   HostManagers
	extraCmds  []func(*base.CmdContext) *cobra.Command
}

type HostManagers struct {
	getMgr HostGetter
	// for update we use simple commands in sake of simplicity
	powerMgr     HostPowerer
	reinstallMgr HostReinstaller
}

// NewCmd builds the slim top-level "hosts" command that only exposes the
// cross-type list. The per-type commands (ebm/kbm/sbm) are top-level commands
// of their own, built by NewEBMCmd/NewKBMCmd/NewSBMCmd.
func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	entitiesMap, err := getHostsEntities()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "Manage hosts",
		Long:  "Manage hosts of different types ( enterprise bare metal, kubernetes baremetal node, scalable baremetal server)",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	// hosts list cmd
	cmd.AddCommand(newListCmd(cmdContext))

	base.AddFormatFlags(cmd)

	return cmd
}

func newHostTypeCmd(cmdContext *base.CmdContext, hostTypeCmd HostTypeCmd) *cobra.Command {
	entitiesMap, err := getHostsEntities()
	if err != nil {
		log.Fatal(err)
	}

	hostCmd := &cobra.Command{
		Use:   hostTypeCmd.use,
		Short: hostTypeCmd.shortDesc,
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	hostCmd.AddCommand(newGetCmd(cmdContext, &hostTypeCmd))

	if hostTypeCmd.managers.powerMgr != nil {
		hostCmd.AddCommand(newPowerCmd(cmdContext, &hostTypeCmd))
		hostCmd.AddCommand(newListPowerFeedsCmd(cmdContext, &hostTypeCmd))
	}
	if hostTypeCmd.managers.reinstallMgr != nil {
		hostCmd.AddCommand(newReinstallCmd(cmdContext, &hostTypeCmd))
	}

	for _, cmdFunc := range hostTypeCmd.extraCmds {
		hostCmd.AddCommand(cmdFunc(cmdContext))
	}

	base.AddFormatFlags(hostCmd)

	return hostCmd
}

func getHostsEntities() (map[string]entities.EntityInterface, error) {
	result := make(map[string]entities.EntityInterface)
	hostsEntity, err := entities.Registry.GetEntityFromValue(serverscom.Host{})
	if err != nil {
		return nil, err
	}
	result["hosts"] = hostsEntity

	ebmEntity, err := entities.Registry.GetEntityFromValue(serverscom.DedicatedServer{})
	if err != nil {
		return nil, err
	}
	result["ebm"] = ebmEntity

	kbmEntity, err := entities.Registry.GetEntityFromValue(serverscom.KubernetesBaremetalNode{})
	if err != nil {
		return nil, err
	}
	result["kbm"] = kbmEntity

	sbmEntity, err := entities.Registry.GetEntityFromValue(serverscom.SBMServer{})
	if err != nil {
		return nil, err
	}
	result["sbm"] = sbmEntity

	hostNetworkEntity, err := entities.Registry.GetEntityFromValue(serverscom.Network{})
	if err != nil {
		return nil, err
	}
	result["list-networks"] = hostNetworkEntity

	return result, nil
}
