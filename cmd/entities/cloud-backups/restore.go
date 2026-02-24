package cloudbackups

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newRestoreCmd(cmdContext *base.CmdContext) *cobra.Command {
	var volumeID string

	cmd := &cobra.Command{
		Use:   "restore <backup-id> --volume-id <volume-id>",
		Short: "Restore a cloud backup",
		Long:  "Restore a cloud backup to a volume",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			required := []string{"volume-id"}
			if err := base.ValidateFlags(cmd, required); err != nil {
				return err
			}

			input := serverscom.CloudBlockStorageBackupRestoreInput{
				VolumeID: volumeID,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			backup, err := scClient.CloudBlockStorageBackups.Restore(ctx, id, input)
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

	cmd.Flags().StringVarP(&volumeID, "volume-id", "", "", "ID of the volume to restore to")

	return cmd
}
