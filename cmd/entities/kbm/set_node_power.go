package kbm

import (
	"context"
	"fmt"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newPowerCmd(cmdContext *base.CmdContext) *cobra.Command {
	var powerFlag string

	cmd := &cobra.Command{
		Use:   "set-node-power <id>",
		Short: "Power on/off/cycle KBM nodes",
		Long:  "Power on/off/cycle KBM nodes by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]

			server, err := powerAction(ctx, scClient, id, powerFlag)
			if err != nil {
				return err
			}

			if server != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(server)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&powerFlag, "power", "", "power command: on|off|cycle (required)")

	return cmd
}

func powerAction(ctx context.Context, client *serverscom.Client, id string, action string) (any, error) {
	actions := map[string]func() (any, error){
		"on": func() (any, error) {
			return client.Hosts.PowerOnKubernetesBaremetalNode(ctx, id)
		},
		"off": func() (any, error) {
			return client.Hosts.PowerOffKubernetesBaremetalNode(ctx, id)
		},
		"cycle": func() (any, error) {
			return client.Hosts.PowerCycleKubernetesBaremetalNode(ctx, id)
		},
	}

	if fn, ok := actions[action]; ok {
		return fn()
	}

	return nil, fmt.Errorf("unsupported power action: %q", action)
}
