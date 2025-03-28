package invoices

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an invoice",
		Long:  "Get an invoice by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			invoice, err := scClient.Invoices.GetBillingInvoice(ctx, id)
			if err != nil {
				return err
			}

			if invoice != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(invoice)
			}
			return nil
		},
	}

	return cmd
}
