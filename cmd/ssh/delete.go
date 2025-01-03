package ssh

import (
	"log"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <fingerprint>",
		Short: "Delete an ssh key",
		Long:  "Delete an ssh key by fingerprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := config.NewManager()
			if err != nil {
				log.Fatal(err)
			}

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := client.NewClient(
				manager.GetToken(),
				manager.GetEndpoint(),
			).SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			fingerprint := args[0]
			return scClient.SSHKeys.Delete(ctx, fingerprint)
		},
	}

	return cmd
}
