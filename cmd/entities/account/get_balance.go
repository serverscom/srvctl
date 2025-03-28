package account

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetBalanceCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Get an account balance info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			balance, err := scClient.Account.GetBalance(ctx)
			if err != nil {
				return err
			}

			if balance != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(balance)
			}
			return nil
		},
	}

	return cmd
}
