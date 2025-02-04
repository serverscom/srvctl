package hosts

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

type HostType struct {
	use        string
	shortDesc  string
	entityName string
	typeFlag   string
	getter     HostGetter
	extraCmds  []func(*base.CmdContext) *cobra.Command
}

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	entitiesMap, err := getHostsEntities()
	if err != nil {
		log.Fatal(err)
	}

	hostTypes := []HostType{
		{
			use:        "ds",
			shortDesc:  "Manage dedicated servers",
			entityName: "Dedicated servers",
			typeFlag:   "dedicated_server",
			getter:     &DedicatedServerGetter{client: cmdContext.GetClient().GetScClient()},
			extraCmds:  []func(*base.CmdContext) *cobra.Command{},
		},
		{
			use:        "kbm",
			shortDesc:  "Manage kubernetes baremetal nodes",
			entityName: "Kubernetes baremetal nodes",
			typeFlag:   "kubernetes_baremetal_node",
			getter:     &KubernetesBaremetalNodeGetter{client: cmdContext.GetClient().GetScClient()},
		},
		{
			use:        "sbm",
			shortDesc:  "Manage scalable baremetal servers",
			entityName: "Scalable baremetal servers",
			typeFlag:   "sbm_server",
			getter:     &SBMServerGetter{client: cmdContext.GetClient().GetScClient()},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// hosts list cmd
	cmd.AddCommand(newListCmd(cmdContext, nil))

	for _, ht := range hostTypes {
		cmd.AddCommand(newHostTypeCmd(cmdContext, ht))
	}

	base.AddFormatFlags(cmd)

	return cmd
}

func newHostTypeCmd(cmdContext *base.CmdContext, hostType HostType) *cobra.Command {
	hostCmd := &cobra.Command{
		Use:   hostType.use,
		Short: hostType.shortDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	hostCmd.AddCommand(newListCmd(cmdContext, &hostType))
	hostCmd.AddCommand(newGetCmd(cmdContext, &hostType))

	for _, cmdFunc := range hostType.extraCmds {
		hostCmd.AddCommand(cmdFunc(cmdContext))
	}

	return hostCmd
}

// func newHostTypeCmd(cmdContext *base.CmdContext, hostType string, shortDesc string, entityName string, typeFlag string) *cobra.Command {
// 	hostCmd := &cobra.Command{
// 		Use:   hostType,
// 		Short: shortDesc,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return cmd.Help()
// 		},
// 	}

// 	listCmd := newListCmd(cmdContext, entityName)
// 	if err := listCmd.Flags().Set("type", typeFlag); err != nil {
// 		log.Fatal(err)
// 	}

// 	hostCmd.AddCommand(listCmd)

// 	switch hostType {
// 	case "ds":
// 		hostCmd.AddCommand(newGetDsCmd(cmdContext))
// 	case "kbm":
// 		hostCmd.AddCommand(newGetKbmCmd(cmdContext))
// 	case "sbm":
// 		hostCmd.AddCommand(newGetSbmCmd(cmdContext))
// 	default:
// 		log.Fatalf("Unknown host type: %s", hostType)
// 	}

// 	return hostCmd
// }

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
