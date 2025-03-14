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
	getMgr       HostGetter
	createMgr    HostCreator
	powerMgr     HostPowerer
	reinstallMgr HostReinstaller
}

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	entitiesMap, err := getHostsEntities()
	if err != nil {
		log.Fatal(err)
	}

	hostTypeCmds := []HostTypeCmd{
		{
			use:        "ds",
			shortDesc:  "Manage dedicated servers",
			entityName: "Dedicated servers",
			typeFlag:   "dedicated_server",
			managers: HostManagers{
				getMgr:       &DSGetMgr{},
				createMgr:    &DSCreateMgr{},
				powerMgr:     &DSPowerMgr{},
				reinstallMgr: &DSReinstallMgr{},
			},
			extraCmds: []func(*base.CmdContext) *cobra.Command{
				newUpdateDSCmd,
				newListDSDriveSlotsCmd,
				newListDSConnectionsCmd,
				newListDSPTRCmd,
				newDSAbortReleaseCmd,
				newDSScheduleReleaseCmd,
			},
		},
		{
			use:        "kbm",
			shortDesc:  "Manage kubernetes baremetal nodes",
			entityName: "Kubernetes baremetal nodes",
			typeFlag:   "kubernetes_baremetal_node",
			managers: HostManagers{
				getMgr:   &KBMGetMgr{},
				powerMgr: &KBMPowerMgr{},
			},
			extraCmds: []func(*base.CmdContext) *cobra.Command{
				newUpdateKBMCmd,
			},
		},
		{
			use:        "sbm",
			shortDesc:  "Manage scalable baremetal servers",
			entityName: "Scalable baremetal servers",
			typeFlag:   "sbm_server",
			managers: HostManagers{
				getMgr:       &SBMGetMgr{},
				createMgr:    &SBMCreateMgr{},
				powerMgr:     &SBMPowerMgr{},
				reinstallMgr: &SBMReinstallMgr{},
			},
			extraCmds: []func(*base.CmdContext) *cobra.Command{
				newUpdateSBMCmd,
				newSBMReleaseCmd,
			},
		},
	}

	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "Manage hosts",
		Long:  "Manage hosts of different types ( dedicated server, kubernetes baremetal node, scalable baremetal server)",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	// hosts list cmd
	cmd.AddCommand(newListCmd(cmdContext, nil))

	for _, ht := range hostTypeCmds {
		cmd.AddCommand(newHostTypeCmd(cmdContext, ht))
	}

	base.AddFormatFlags(cmd)

	return cmd
}

func newHostTypeCmd(cmdContext *base.CmdContext, hostTypeCmd HostTypeCmd) *cobra.Command {
	hostCmd := &cobra.Command{
		Use:   hostTypeCmd.use,
		Short: hostTypeCmd.shortDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	hostCmd.AddCommand(newListCmd(cmdContext, &hostTypeCmd))
	hostCmd.AddCommand(newGetCmd(cmdContext, &hostTypeCmd))

	if hostTypeCmd.managers.getMgr != nil {
		hostCmd.AddCommand(newAddCmd(cmdContext, &hostTypeCmd))
	}
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

	return hostCmd
}

func getHostsEntities() (map[string]entities.EntityInterface, error) {
	result := make(map[string]entities.EntityInterface)
	hostsEntity, err := entities.Registry.GetEntityFromValue(serverscom.Host{})
	if err != nil {
		return nil, err
	}
	result["hosts"] = hostsEntity

	dsEntity, err := entities.Registry.GetEntityFromValue(serverscom.DedicatedServer{})
	if err != nil {
		return nil, err
	}
	result["ds"] = dsEntity

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

	return result, nil
}
