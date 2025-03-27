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
	getMgr    LBGetter
	createMgr LBCreator
	updateMgr LBUpdater
	deleteMgr LBDeleter
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
			managers: LBManagers{
				getMgr:    &LBL4GetMgr{},
				createMgr: &LBL4CreateMgr{},
				updateMgr: &LBL4UpdateMgr{},
				deleteMgr: &LBL4DeleteMgr{},
			},
		},
		{
			use:        "l7",
			shortDesc:  "Manage L7 load balancers",
			entityName: "L7 load balancers",
			typeFlag:   "l7",
			managers: LBManagers{
				getMgr:    &LBL7GetMgr{},
				createMgr: &LBL7CreateMgr{},
				updateMgr: &LBL7UpdateMgr{},
				deleteMgr: &LBL7DeleteMgr{},
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

func newLBTypeCmd(cmdContext *base.CmdContext, lbType LBTypeCmd) *cobra.Command {
	LBCmd := &cobra.Command{
		Use:   lbType.use,
		Short: lbType.shortDesc,
		Args:  base.NoArgs,
		Run:   base.UsageRun,
	}

	LBCmd.AddCommand(newListCmd(cmdContext, &lbType))
	LBCmd.AddCommand(newGetCmd(cmdContext, &lbType))

	if lbType.managers.createMgr != nil {
		LBCmd.AddCommand(newAddCmd(cmdContext, &lbType))
	}
	if lbType.managers.updateMgr != nil {
		LBCmd.AddCommand(newUpdateCmd(cmdContext, &lbType))
	}
	if lbType.managers.deleteMgr != nil {
		LBCmd.AddCommand(newDeleteCmd(cmdContext, &lbType))
	}

	for _, cmdFunc := range lbType.extraCmds {
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
