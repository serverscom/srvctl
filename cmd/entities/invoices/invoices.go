package invoices

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	invoiceEntity, err := entities.Registry.GetEntityFromValue(serverscom.Invoice{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["invoices"] = invoiceEntity
	cmd := &cobra.Command{
		Use:   "invoices",
		Short: "Manage invoices",
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
