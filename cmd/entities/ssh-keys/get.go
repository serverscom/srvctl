package sshkeys

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <fingerprint>",
		Short: "Get an ssh key",
		Long:  "Get an ssh key by fingerprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			fingerprint := args[0]
			sshKey, err := scClient.SSHKeys.Get(ctx, fingerprint)
			if err != nil {
				return err
			}

			if sshKey != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(sshKey)
			}
			return nil
		},
	}

	return cmd
}
