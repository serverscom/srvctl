package ssh

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var name string
	var labels []string

	cmd := &cobra.Command{
		Use:   "update <fingerprint>",
		Short: "Update an ssh key",
		Long:  "Update an ssh key by fingerprint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap, err := base.ParseLabels(labels)
			if err != nil {
				log.Fatal(err)
			}
			input := serverscom.SSHKeyUpdateInput{
				Name:   name,
				Labels: labelsMap,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			fingerprint := args[0]
			sshKey, err := scClient.SSHKeys.Update(ctx, fingerprint, input)
			if err != nil {
				return err
			}

			if sshKey != nil {
				outputFormat, _ := manager.GetResolvedStringValue(cmd, "output")
				formatter := output.NewFormatter(cmd.OutOrStdout())
				return formatter.Format([]serverscom.SSHKey{*sshKey}, outputFormat)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "string")
	cmd.Flags().StringArrayVarP(&labels, "labels", "l", []string{}, "string in JSON format")

	return cmd
}
