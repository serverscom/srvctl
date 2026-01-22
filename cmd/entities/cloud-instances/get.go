package cloudinstances

import (
	"log"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newGetCmd(cmdContext *base.CmdContext) *cobra.Command {
	var instanceID string

	cmd := &cobra.Command{
		Use:   "get --instance-id <id>",
		Short: "Get cloud instance by instance id",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			item, err := scClient.CloudComputingInstances.Get(ctx, instanceID)
			if err != nil {
				return err
			}

			formatter := cmdContext.GetOrCreateFormatter(cmd)
			return formatter.Format(item)
		},
	}

	cmd.Flags().StringVar(&instanceID, "instance-id", "", "Cloud instance ID")
	if err := cmd.MarkFlagRequired("instance-id"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
