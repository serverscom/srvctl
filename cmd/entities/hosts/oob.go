package hosts

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetDSOOBCredsCmd(cmdContext *base.CmdContext) *cobra.Command {
	var fingerprint string

	cmd := &cobra.Command{
		Use:   "get-oob-credentials <id>",
		Short: ("Get a dedicated server oob credentials"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := scClient.Hosts.GetDedicatedServerOOBCredentials(ctx, id, map[string]string{"fingerprint": fingerprint})
			if err != nil {
				return err
			}

			if entity != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(entity)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&fingerprint, "fingerprint", "", "GPG key fingerprint (required)")
	_ = cmd.MarkFlagRequired("fingerprint")

	return cmd
}
