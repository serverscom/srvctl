package cloudinstances

import (
	"fmt"
	"log"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newRescueModeCmd(cmdContext *base.CmdContext) *cobra.Command {
	var command string

	cmd := &cobra.Command{
		Use:   "rescue-mode <instance-id>",
		Short: "Activate or deactivate rescue mode of cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			var out any
			var err error
			switch command {
			case "activate":
				out, err = scClient.CloudComputingInstances.Rescue(ctx, instanceID)
			case "deactivate":
				out, err = scClient.CloudComputingInstances.Unrescue(ctx, instanceID)
			default:
				return fmt.Errorf("invalid command: %s", command)
			}
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

	cmd.Flags().StringVar(&command, "command", "", "Command: activate or deactivate")
	if err := cmd.MarkFlagRequired("command"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
