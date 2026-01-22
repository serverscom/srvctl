package cloudinstances

import (
	"fmt"
	"log"

	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newPowerCmd(cmdContext *base.CmdContext) *cobra.Command {
	var command string

	cmd := &cobra.Command{
		Use:   "power <instance-id>",
		Short: "Power on or off cloud instance by instance ID",
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
			case "on":
				out, err = scClient.CloudComputingInstances.PowerOn(ctx, instanceID)
			case "off":
				out, err = scClient.CloudComputingInstances.PowerOff(ctx, instanceID)
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

	cmd.Flags().StringVar(&command, "command", "", "Command: on or off")
	if err := cmd.MarkFlagRequired("command"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
