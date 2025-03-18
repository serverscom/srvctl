package ssl

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <fingerprint>",
		Short: "Get an ssl certificate key",
		Long:  "Get an ssl certificate by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			fingerprint := args[0]
			sslCert, err := scClient.SSLCertificates.GetCustom(ctx, fingerprint)
			if err != nil {
				return err
			}

			if sslCert != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(sslCert)
			}
			return nil
		},
	}

	return cmd
}
