package ssh

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <fingerprint>",
		Short: "Get an ssh key",
		Long:  "Get an ssh key by fingerprint",
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
			sshKey, err := scClient.SSHKeys.Get(ctx, fingerprint)
			if err != nil {
				return err
			}

			if sshKey != nil {
				outputFormat, _ := manager.GetResolvedStringValue(cmd, "output")
				return output.Format([]serverscom.SSHKey{*sshKey}, outputFormat)
			}
			return nil
		},
	}

	return cmd
}
