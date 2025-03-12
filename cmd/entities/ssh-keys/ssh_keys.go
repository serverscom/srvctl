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
	entitiesMap := make(map[string]entities.EntityInterface)
	entitiesMap["ssh-keys"] = sshEntity
	cmd := &cobra.Command{
		Use:   "ssh-keys",
		Short: "Manage ssh keys",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newAddCmd(cmdContext),
		newGetCmd(cmdContext),
		newUpdateCmd(cmdContext),
		newDeleteCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}
