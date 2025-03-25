package loadbalancers

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

type LBTypeCmd struct {
	use        string
	shortDesc  string
	entityName string
	typeFlag   string
	managers   LBManagers
	extraCmds  []func(*base.CmdContext) *cobra.Command
}

type LBManagers struct {
	// getMgr    LBGetter
	// createMgr LBCreator
	// for update we use simple commands in sake of simplicity
	// deleteMgr LBDeleter
}

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	entitiesMap, err := getLBEntities()
	if err != nil {
		log.Fatal(err)
	}

	LBTypeCmds := []LBTypeCmd{
		{
			use:        "l4",
			shortDesc:  "Manage L4 load balancers",
			entityName: "L4 load balancers",
			typeFlag:   "l4",
			managers:   LBManagers{
				// getMgr:    &LBCustomGetMgr{},
				// createMgr: &LBCustomCreateMgr{},
				// deleteMgr: &LBCustomDeleteMgr{},
			},
			extraCmds: []func(*base.CmdContext) *cobra.Command{
				// newUpdateCustomCmd,
			},
		},
		{
			use:        "l7",
			shortDesc:  "Manage L7 load balancers",
			entityName: "L7 load balancers",
			typeFlag:   "l7",
			managers:   LBManagers{
				// getMgr:    &LBLeGetMgr{},
				// deleteMgr: &LBLeDeleteMgr{},
			},
			extraCmds: []func(*base.CmdContext) *cobra.Command{
				// newUpdateLeCmd,
			},
		},
	}

	cmd := &cobra.Command{
		Use:   "lb",
		Short: "Manage load balancers",
		Long:  "Manage load balancers of different types ( l4, l7 )",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	// lb list cmd
	cmd.AddCommand(newListCmd(cmdContext, nil))

	for _, st := range LBTypeCmds {
		cmd.AddCommand(newLBTypeCmd(cmdContext, st))
	}

	base.AddFormatFlags(cmd)

	return cmd
}

func newLBTypeCmd(cmdContext *base.CmdContext, LBTypeCmd LBTypeCmd) *cobra.Command {
	LBCmd := &cobra.Command{
		Use:   LBTypeCmd.use,
		Short: LBTypeCmd.shortDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	LBCmd.AddCommand(newListCmd(cmdContext, &LBTypeCmd))
	// LBCmd.AddCommand(newGetCmd(cmdContext, &LBTypeCmd))

	// if LBTypeCmd.managers.createMgr != nil {
	// LBCmd.AddCommand(newAddCmd(cmdContext, &LBTypeCmd))
	// }
	// if LBTypeCmd.managers.deleteMgr != nil {
	// LBCmd.AddCommand(newDeleteCmd(cmdContext, &LBTypeCmd))
	// }

	for _, cmdFunc := range LBTypeCmd.extraCmds {
		LBCmd.AddCommand(cmdFunc(cmdContext))
	}

	return LBCmd
}

func getLBEntities() (map[string]entities.EntityInterface, error) {
	result := make(map[string]entities.EntityInterface)
	lbEntity, err := entities.Registry.GetEntityFromValue(serverscom.LoadBalancer{})
	if err != nil {
		return nil, err
	}
	result["lb"] = lbEntity
	result["list"] = lbEntity

	l4Entity, err := entities.Registry.GetEntityFromValue(serverscom.L4LoadBalancer{})
	if err != nil {
		return nil, err
	}
	result["l4"] = l4Entity

	l7Entity, err := entities.Registry.GetEntityFromValue(serverscom.L7LoadBalancer{})
	if err != nil {
		return nil, err
	}
	result["l7"] = l7Entity

	return result, nil
}
