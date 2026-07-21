package dns

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newDeleteRecordCmd(cmdContext *base.CmdContext) *cobra.Command {
	var recordID string

	cmd := &cobra.Command{
		Use:   "delete-record <domain-id> --record-id <record-id>",
		Short: "Delete a DNS record",
		Long:  "Delete a DNS record by id within a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := base.ValidateFlags(cmd, []string{"record-id"}); err != nil {
				return err
			}

			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			domainID := args[0]
			return scClient.DNS.DeleteRecord(ctx, domainID, recordID)
		},
	}

	cmd.Flags().StringVarP(&recordID, "record-id", "", "", "A DNS record id")

	return cmd
}
