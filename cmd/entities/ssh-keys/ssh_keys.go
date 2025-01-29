package sshkeys

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	sshEntity, err := entities.Registry.GetEntityFromValue(serverscom.SSHKey{})
	if err != nil {
		log.Fatal(err)
	}
	cmd := &cobra.Command{
		Use:   "ssh-keys",
		Short: "Manage ssh keys",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, sshEntity),
			base.CheckEmptyContexts(cmdContext),
		),
		// empty RunE to support flags for ssh-keys command itself
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newAddCmd(cmdContext),
		newGetCmd(cmdContext),
		newUpdateCmd(cmdContext),
		newDeleteCmd(cmdContext),
	)

	cmd.PersistentFlags().StringArrayP("field", "f", []string{}, "output only these fields, can be specified multiple times")
	cmd.PersistentFlags().Bool("field-list", false, "list available fields")
	cmd.PersistentFlags().Bool("page-view", false, "use page view format")
	cmd.PersistentFlags().StringP("template", "t", "", "go template string to output in specified format")

	return cmd
}
