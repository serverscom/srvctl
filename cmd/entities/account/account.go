package account

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	accountBalanceEntity, err := entities.Registry.GetEntityFromValue(serverscom.AccountBalance{})
	if err != nil {
		log.Fatal(err)
	}
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["account"] = accountBalanceEntity
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage account operations",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newGetBalanceCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
