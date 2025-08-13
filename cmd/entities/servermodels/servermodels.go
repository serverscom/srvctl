package servermodels

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	entitiesMap, err := getServerModelOptionsEntities()
	if err != nil {
		log.Fatal(err)
	}
	cmd := &cobra.Command{
		Use:   "server-models",
		Short: "Manage server models",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newGetCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}

func getServerModelOptionsEntities() (map[string]entities.EntityInterface, error) {
	result := make(map[string]entities.EntityInterface)
	getEntity, err := entities.Registry.GetEntityFromValue(serverscom.ServerModelOptionDetail{})
	if err != nil {
		return nil, err
	}
	listEntity, err := entities.Registry.GetEntityFromValue(serverscom.ServerModelOption{})
	if err != nil {
		return nil, err
	}

	result["get"] = getEntity
	result["list"] = listEntity

	return result, nil
}
