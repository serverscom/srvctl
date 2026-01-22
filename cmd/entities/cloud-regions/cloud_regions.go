package cloudregions

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	cloudEntity, err := entities.Registry.GetEntityFromValue(serverscom.CloudComputingRegion{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["cloud-regions"] = cloudEntity

	cmd := &cobra.Command{
		Use:   "cloud-regions",
		Short: "Manage cloud regions",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newGetCredentialsCmd(cmdContext),
		newListImagesCmd(cmdContext),
		newListFlavorsCmd(cmdContext),
		newListSnapshotsCmd(cmdContext),
		newAddSnapshotCmd(cmdContext),
		newDeleteSnapshotCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
