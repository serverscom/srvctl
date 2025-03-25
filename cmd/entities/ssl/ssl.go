package ssl

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

type SSLTypeCmd struct {
	use        string
	shortDesc  string
	entityName string
	typeFlag   string
	managers   SSLManagers
	extraCmds  []func(*base.CmdContext) *cobra.Command
}

type SSLManagers struct {
	getMgr    SSLGetter
	createMgr SSLCreator
	// for update we use simple commands in sake of simplicity
	deleteMgr SSLDeleter
}

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	entitiesMap, err := getSSLEntities()
	if err != nil {
		log.Fatal(err)
	}

	sslTypeCmds := []SSLTypeCmd{
		{
			use:        "custom",
			shortDesc:  "Manage ssl custom certificates",
			entityName: "Custom SSL certificates",
			typeFlag:   "custom",
			managers: SSLManagers{
				getMgr:    &SSLCustomGetMgr{},
				createMgr: &SSLCustomCreateMgr{},
				deleteMgr: &SSLCustomDeleteMgr{},
			},
			extraCmds: []func(*base.CmdContext) *cobra.Command{
				newUpdateCustomCmd,
			},
		},
		{
			use:        "le",
			shortDesc:  "Manage ssl le certificates",
			entityName: "LE SSL certificates",
			typeFlag:   "le",
			managers: SSLManagers{
				getMgr:    &SSLLeGetMgr{},
				deleteMgr: &SSLLeDeleteMgr{},
			},
			extraCmds: []func(*base.CmdContext) *cobra.Command{
				newUpdateLeCmd,
			},
		},
	}

	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "Manage SSL certificates",
		Long:  "Manage SSL certificates of different types ( custom, let's encrypt )",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	// ssl list cmd
	cmd.AddCommand(newListCmd(cmdContext, nil))

	for _, st := range sslTypeCmds {
		cmd.AddCommand(newSSLTypeCmd(cmdContext, st))
	}

	base.AddFormatFlags(cmd)

	return cmd
}

func newSSLTypeCmd(cmdContext *base.CmdContext, sslTypeCmd SSLTypeCmd) *cobra.Command {
	sslCmd := &cobra.Command{
		Use:   sslTypeCmd.use,
		Short: sslTypeCmd.shortDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	sslCmd.AddCommand(newListCmd(cmdContext, &sslTypeCmd))
	sslCmd.AddCommand(newGetCmd(cmdContext, &sslTypeCmd))

	if sslTypeCmd.managers.createMgr != nil {
		sslCmd.AddCommand(newAddCmd(cmdContext, &sslTypeCmd))
	}
	if sslTypeCmd.managers.deleteMgr != nil {
		sslCmd.AddCommand(newDeleteCmd(cmdContext, &sslTypeCmd))
	}

	for _, cmdFunc := range sslTypeCmd.extraCmds {
		sslCmd.AddCommand(cmdFunc(cmdContext))
	}

	return sslCmd
}

func getSSLEntities() (map[string]entities.EntityInterface, error) {
	result := make(map[string]entities.EntityInterface)
	sslEntity, err := entities.Registry.GetEntityFromValue(serverscom.SSLCertificate{})
	if err != nil {
		return nil, err
	}
	result["ssl"] = sslEntity
	result["list"] = sslEntity

	sslCustomEntity, err := entities.Registry.GetEntityFromValue(serverscom.SSLCertificateCustom{})
	if err != nil {
		return nil, err
	}
	result["custom"] = sslCustomEntity

	sslLeEntity, err := entities.Registry.GetEntityFromValue(serverscom.SSLCertificateLE{})
	if err != nil {
		return nil, err
	}
	result["le"] = sslLeEntity

	return result, nil
}
