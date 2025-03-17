package ssl

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	sslCertEntity, err := entities.Registry.GetEntityFromValue(serverscom.SSLCertificate{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["ssl"] = sslCertEntity
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "Manage SSL certificates",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		// newAddCmd(cmdContext),
		// newGetCmd(cmdContext),
		// newUpdateCmd(cmdContext),
		// newDeleteCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
