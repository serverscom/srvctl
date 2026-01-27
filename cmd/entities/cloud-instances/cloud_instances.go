package cloudinstances

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	cloudEntity, err := entities.Registry.GetEntityFromValue(serverscom.CloudComputingInstance{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["cloud-instances"] = cloudEntity

	cmd := &cobra.Command{
		Use:   "cloud-instances",
		Short: "Manage cloud instances",
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
		newAddCmd(cmdContext),
		newUpdateCmd(cmdContext),
		newDeleteCmd(cmdContext),
		newReinstallCmd(cmdContext),
		newUpgradeCmd(cmdContext),
		newUpgradeRevertCmd(cmdContext),
		newUpgradeApproveCmd(cmdContext),
		newRebootCmd(cmdContext),
		newRescueModeCmd(cmdContext),
		newPowerCmd(cmdContext),
		newListPTRCmd(cmdContext),
		newAddPTRCmd(cmdContext),
		newDeletePTRCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
