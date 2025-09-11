package kbm

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	kbmEntity, err := entities.Registry.GetEntityFromValue(serverscom.KubernetesBaremetalNode{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["kbm"] = kbmEntity
	cmd := &cobra.Command{
		Use:   "kbm",
		Short: "Manage kubernetes baremetal nodes",
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
		newGetPowerFeedsCmd(cmdContext),
		newUpdateCmd(cmdContext),
		newPowerCmd(cmdContext),
		newListNetworksCmd(cmdContext),
		newListSlotsCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
