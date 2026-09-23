package dns

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type UpdateRecordFlags struct {
	RecordID string
	Name     string
	Data     string
	TTL      int
	Priority int
}

func newUpdateRecordCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &UpdateRecordFlags{}

	cmd := &cobra.Command{
		Use:   "update-record <domain-id> --record-id <record-id>",
		Short: "Update a DNS record",
		Long:  "Update a DNS record by id within a domain. Record type cannot be changed; SOA/PTR records are not editable",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := base.ValidateFlags(cmd, []string{"record-id"}); err != nil {
				return err
			}

			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.DNSRecordUpdateInput{}
			flags.FillInput(cmd, &input)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			domainID := args[0]
			record, err := scClient.DNS.UpdateRecord(ctx, domainID, flags.RecordID, input)
			if err != nil {
				return err
			}

			if record != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(record)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.RecordID, "record-id", "", "", "A DNS record id")
	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "A name of a DNS record")
	cmd.Flags().StringVarP(&flags.Data, "data", "", "", "A data of a DNS record")
	cmd.Flags().IntVarP(&flags.TTL, "ttl", "", 0, "A TTL of a DNS record in seconds")
	cmd.Flags().IntVarP(&flags.Priority, "priority", "", 0, "A priority of a DNS record (required by convention for MX/SRV, ignored for other types)")

	return cmd
}

func (f *UpdateRecordFlags) FillInput(cmd *cobra.Command, input *serverscom.DNSRecordUpdateInput) {
	if cmd.Flags().Changed("name") {
		input.Name = f.Name
	}
	if cmd.Flags().Changed("data") {
		input.Data = f.Data
	}
	if cmd.Flags().Changed("ttl") {
		input.TTL = &f.TTL
	}
	if cmd.Flags().Changed("priority") {
		input.Priority = &f.Priority
	}
}
