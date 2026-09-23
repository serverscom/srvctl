package dns

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <domain-id>",
		Short: "Get a DNS domain",
		Long:  "Get a DNS domain by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			domainID := args[0]
			domain, err := scClient.DNS.GetDomain(ctx, domainID)
			if err != nil {
				return err
			}

			if domain != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(domain)
			}
			return nil
		},
	}

	return cmd
}
