package ssh

import (
	"encoding/json"
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var name, labels string

	cmd := &cobra.Command{
		Use:   "update <fingerprint>",
		Short: "Update an ssh key",
		Long:  "Update an ssh key by fingerprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := config.NewManager()
			if err != nil {
				log.Fatal(err)
			}

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap := make(map[string]string)
			if labels != "" {
				if err := json.Unmarshal([]byte(labels), &labelsMap); err != nil {
					log.Fatal(err)
				}
			}
			input := serverscom.SSHKeyUpdateInput{
				Name:   name,
				Labels: labelsMap,
			}

			scClient := client.NewClient(
				manager.GetToken(),
				manager.GetEndpoint(),
			).SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			fingerprint := args[0]
			sshKey, err := scClient.SSHKeys.Update(ctx, fingerprint, input)
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

	cmd.Flags().StringVarP(&name, "name", "n", "", "string")
	cmd.Flags().StringVarP(&labels, "labels", "l", "", "string in JSON format")

	return cmd
}
