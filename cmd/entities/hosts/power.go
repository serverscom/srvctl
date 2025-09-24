package hosts

import (
	"context"
	"fmt"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type HostPowerer interface {
	PowerAction(ctx context.Context, client *serverscom.Client, id string, action string) (any, error)
	ListPowerFeeds(ctx context.Context, client *serverscom.Client, id string) (any, error)
}

type DSPowerMgr struct{}

func (p *DSPowerMgr) PowerAction(ctx context.Context, client *serverscom.Client, id string, action string) (any, error) {
	switch action {
	case "on":
		return client.Hosts.PowerOnDedicatedServer(ctx, id)
	case "off":
		return client.Hosts.PowerOffDedicatedServer(ctx, id)
	case "cycle":
		return client.Hosts.PowerCycleDedicatedServer(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported power action: %s", action)
	}
}

func (p *DSPowerMgr) ListPowerFeeds(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.Hosts.DedicatedServerPowerFeeds(ctx, id)
}

type KBMPowerMgr struct{}

func (p *KBMPowerMgr) PowerAction(ctx context.Context, client *serverscom.Client, id string, action string) (any, error) {
	switch action {
	case "on":
		return client.Hosts.PowerOnKubernetesBaremetalNode(ctx, id)
	case "off":
		return client.Hosts.PowerOffKubernetesBaremetalNode(ctx, id)
	case "cycle":
		return client.Hosts.PowerCycleKubernetesBaremetalNode(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported power action: %s", action)
	}
}

func (p *KBMPowerMgr) ListPowerFeeds(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.Hosts.KubernetesBaremetalNodePowerFeeds(ctx, id)
}

type SBMPowerMgr struct{}

func (p *SBMPowerMgr) PowerAction(ctx context.Context, client *serverscom.Client, id string, action string) (any, error) {
	switch action {
	case "on":
		return client.Hosts.PowerOnSBMServer(ctx, id)
	case "off":
		return client.Hosts.PowerOffSBMServer(ctx, id)
	case "cycle":
		return client.Hosts.PowerCycleSBMServer(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported power action: %s", action)
	}
}

func (p *SBMPowerMgr) ListPowerFeeds(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	// return client.Hosts.TODO(ctx, id)
	return nil, nil
}

func newPowerCmd(cmdContext *base.CmdContext, hostType *HostTypeCmd) *cobra.Command {
	var commandFlag string

	cmd := &cobra.Command{
		Use:   "power <id>",
		Short: fmt.Sprintf("Send power command for %s", hostType.entityName),
		Long:  fmt.Sprintf("Send power command for %s by id", hostType.entityName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			server, err := hostType.managers.powerMgr.PowerAction(ctx, scClient, id, commandFlag)
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

	cmd.Flags().StringVar(&commandFlag, "command", "", "power command")

	return cmd
}

func newListPowerFeedsCmd(cmdContext *base.CmdContext, hostType *HostTypeCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-power-feeds <id>",
		Short: fmt.Sprintf("List power feeds for a %s", hostType.entityName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			server, err := hostType.managers.powerMgr.ListPowerFeeds(ctx, scClient, id)
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

	return cmd
}
