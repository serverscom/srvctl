package dns

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newDelegationDataCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegation-data",
		Short: "Get DNS delegation data",
		Long:  "Get the nameservers and required TXT record to configure at the domain registrar to delegate/verify a DNS zone",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			data, err := scClient.DNS.GetDelegationData(ctx)
			if err != nil {
				return err
			}

			if data != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(data)
			}
			return nil
		},
	}

	return cmd
}
