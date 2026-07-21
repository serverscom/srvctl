package dns

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	dnsDomainEntity, err := entities.Registry.GetEntityFromValue(serverscom.DNSDomain{})
	if err != nil {
		log.Fatal(err)
	}
	dnsRecordEntity, err := entities.Registry.GetEntityFromValue(serverscom.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}
	dnsDelegationEntity, err := entities.Registry.GetEntityFromValue(serverscom.DNSDomainDelegationData{})
	if err != nil {
		log.Fatal(err)
	}

	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["dns"] = dnsDomainEntity
	entitiesMap["list-records"] = dnsRecordEntity
	entitiesMap["get-record"] = dnsRecordEntity
	entitiesMap["add-record"] = dnsRecordEntity
	entitiesMap["update-record"] = dnsRecordEntity
	entitiesMap["delegation-data"] = dnsDelegationEntity

	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Manage DNS domains and records",
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
		newDelegationDataCmd(cmdContext),
		newListRecordsCmd(cmdContext),
		newGetRecordCmd(cmdContext),
		newAddRecordCmd(cmdContext),
		newUpdateRecordCmd(cmdContext),
		newDeleteRecordCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
