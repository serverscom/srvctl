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

func newAddCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: "Add an ssh key",
		Long:  "Add a new SSH key to account",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := config.NewManager()
			if err != nil {
				log.Fatal(err)
			}

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.SSHKeyCreateInput{}
			if err := base.LoadInputFromFile(path, input); err != nil {
				log.Fatal(err)
			}

			scClient := client.NewClient(
				manager.GetToken(),
				manager.GetEndpoint(),
			).SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			sshKey, err := scClient.SSHKeys.Create(ctx, *input)
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

	cmd.Flags().StringVarP(&path, "input", "i", "", "/path/to/create-file.json")
	if err := cmd.MarkFlagRequired("input"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
