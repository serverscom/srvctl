package rbsvolumes

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	rbsVolumeEntity, err := entities.Registry.GetEntityFromValue(serverscom.RemoteBlockStorageVolume{})
	if err != nil {
		log.Fatal(err)
	}
	rbsCredentialsEntity, err := entities.Registry.GetEntityFromValue(serverscom.RemoteBlockStorageVolumeCredentials{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["rbs"] = rbsVolumeEntity
	entitiesMap["rbs-credentials"] = rbsCredentialsEntity
	cmd := &cobra.Command{
		Use:   "rbs",
		Short: "Manage RBS servers",
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
		newGetCredentialsCmd(cmdContext),
		newResetCredentialsCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
