package sshkeys

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: "Add an ssh key",
		Long:  "Add a new SSH key to account",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.SSHKeyCreateInput{}
			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			sshKey, err := scClient.SSHKeys.Create(ctx, *input)
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

	cmd.Flags().StringVarP(&path, "input", "i", "", "path to input file or '-' to read from stdin")
	if err := cmd.MarkFlagRequired("input"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
