package networkpools

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	networkPoolEntity, err := entities.Registry.GetEntityFromValue(serverscom.NetworkPool{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["network-pools"] = networkPoolEntity
	cmd := &cobra.Command{
		Use:   "network-pools",
		Short: "Manage network pools",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newListSubnetsCmd(cmdContext),
		newGetCmd(cmdContext),
		newGetSubnetCmd(cmdContext),
		newCreateSubnetCmd(cmdContext),
		newUpdateCmd(cmdContext),
		newUpdateSubnetCmd(cmdContext),
		newDeleteSubnetCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
