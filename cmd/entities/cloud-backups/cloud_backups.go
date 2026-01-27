package cloudbackups

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	backupEntity, err := entities.Registry.GetEntityFromValue(serverscom.CloudBlockStorageBackup{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["cloud-backups"] = backupEntity
	cmd := &cobra.Command{
		Use:   "cloud-backups",
		Short: "Manage cloud backups",
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
		newRestoreCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
