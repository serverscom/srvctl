package dns

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DNSDomain] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.DNS.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.DNSDomain]{},
		&base.SearchPatternOption[serverscom.DNSDomain]{},
		&base.DelegationStatusOption[serverscom.DNSDomain]{},
		&base.LabelSelectorOption[serverscom.DNSDomain]{},
	)

	return base.NewListCmd("list", "DNS domains", factory, cmdContext, opts...)
}

func newListRecordsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.DNSRecord] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.DNS.Records(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.DNSRecord]{},
		&base.DNSRecordTypeOption[serverscom.DNSRecord]{},
	)

	cmd := base.NewListCmd("list-records", "DNS records", factory, cmdContext, opts...)
	cmd.Use = "list-records <domain-id>"
	cmd.Args = cobra.ExactArgs(1)

	return cmd
}
