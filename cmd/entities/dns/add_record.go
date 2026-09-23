package dns

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddRecordFlags struct {
	Skeleton  bool
	InputPath string
	Type      string
	Name      string
	Data      string
	TTL       int
	Priority  int
}

func newAddRecordCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &AddRecordFlags{}

	cmd := &cobra.Command{
		Use:   "add-record <domain-id>",
		Short: "Add a DNS record",
		Long:  "Add a new DNS record to a domain",
		Args:  base.SkeletonOrExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)

			if flags.Skeleton {
				return formatter.FormatSkeleton("dns/add-record.json")
			}

			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.DNSRecordCreateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else {
				required := []string{"type", "name", "data"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			flags.FillInput(cmd, input)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			domainID := args[0]
			record, err := scClient.DNS.CreateRecord(ctx, domainID, *input)
			if err != nil {
				return err
			}

			if record != nil {
				return formatter.Format(record)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	cmd.Flags().StringVarP(&flags.Type, "type", "", "", "A type of a DNS record (A, AAAA, CNAME, MX, TXT, NS, SRV, CAA, ALIAS)")
	cmd.Flags().StringVarP(&flags.Name, "name", "n", "", "A name of a DNS record")
	cmd.Flags().StringVarP(&flags.Data, "data", "", "", "A data of a DNS record")
	cmd.Flags().IntVarP(&flags.TTL, "ttl", "", 0, "A TTL of a DNS record in seconds")
	cmd.Flags().IntVarP(&flags.Priority, "priority", "", 0, "A priority of a DNS record (required by convention for MX/SRV, ignored for other types)")

	return cmd
}

func (f *AddRecordFlags) FillInput(cmd *cobra.Command, input *serverscom.DNSRecordCreateInput) {
	if cmd.Flags().Changed("type") {
		input.Type = serverscom.DNSRecordType(f.Type)
	}
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
