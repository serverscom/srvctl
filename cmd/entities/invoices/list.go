package invoices

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.InvoiceList] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.Invoices.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.InvoiceList]{},
		&base.StatusOption[serverscom.InvoiceList]{},
		&base.InvoiceTypeOption[serverscom.InvoiceList]{},
		&base.ParentIDOption[serverscom.InvoiceList]{},
		&base.CurrencyOption[serverscom.InvoiceList]{},
		&base.StartDateOption[serverscom.InvoiceList]{},
		&base.EndDateOption[serverscom.InvoiceList]{},
	)

	return base.NewListCmd("list", "invoices", factory, cmdContext, opts...)
}
