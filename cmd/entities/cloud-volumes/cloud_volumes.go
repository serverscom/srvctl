package cloudvolumes

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	cloudVolumeEntity, err := entities.Registry.GetEntityFromValue(serverscom.CloudBlockStorageVolume{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["cloud-volumes"] = cloudVolumeEntity
	cmd := &cobra.Command{
		Use:   "cloud-volumes",
		Short: "Manage cloud volumes",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newAddCmd(cmdContext),
		newGetCmd(cmdContext),
		newUpdateCmd(cmdContext),
		newDeleteCmd(cmdContext),
		newVolumeAttachCmd(cmdContext),
		newVolumeDetachCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
