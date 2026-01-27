package cloudinstances

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newReinstallCmd(cmdContext *base.CmdContext) *cobra.Command {
	var imageID string
	var userData string

	cmd := &cobra.Command{
		Use:   "reinstall <instance-id>",
		Short: "Reinstall cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := serverscom.CloudComputingInstanceReinstallInput{
				ImageID:  imageID,
				UserData: userData,
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			out, err := scClient.CloudComputingInstances.Reinstall(ctx, instanceID, input)
			if err != nil {
				return err
			}

			if out != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(out)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&imageID, "image-id", "", "Image ID for reinstall")
	if err := cmd.MarkFlagRequired("image-id"); err != nil {
		log.Fatal(err)
	}
	cmd.Flags().StringVar(&userData, "user-data", "", "User data for reinstall")

	return cmd
}
