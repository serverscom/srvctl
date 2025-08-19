package l2segments

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	l2Entity, err := entities.Registry.GetEntityFromValue(serverscom.L2Segment{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["l2"] = l2Entity
	cmd := &cobra.Command{
		Use:   "l2",
		Short: "Manage L2 segments",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newListGroupsCmd(cmdContext),
		newListMembersCmd(cmdContext),
		newListNetworksCmd(cmdContext),
		newGetCmd(cmdContext),
		newAddCmd(cmdContext),
		newUpdateL2Cmd(cmdContext),
		newUpdateL2NetworksCmd(cmdContext),
		newDeleteCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
