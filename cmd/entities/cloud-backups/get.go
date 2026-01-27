package cloudbackups

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <backup-id>",
		Short: "Get a cloud backup",
		Long:  "Get a cloud backup by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			backup, err := scClient.CloudBlockStorageBackups.Get(ctx, id)
			if err != nil {
				return err
			}

			if backup != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(backup)
			}
			return nil
		},
	}

	return cmd
}
