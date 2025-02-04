package hosts

import (
	"context"
	"fmt"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type HostGetter interface {
	Get(ctx context.Context, client *serverscom.Client, id string) (interface{}, error)
}

type DedicatedServerGetter struct{}

func (g *DedicatedServerGetter) Get(ctx context.Context, client *serverscom.Client, id string) (interface{}, error) {
	return client.Hosts.GetDedicatedServer(ctx, id)
}

type KubernetesBaremetalNodeGetter struct{}

func (g *KubernetesBaremetalNodeGetter) Get(ctx context.Context, client *serverscom.Client, id string) (interface{}, error) {
	return client.Hosts.GetKubernetesBaremetalNode(ctx, id)
}

type SBMServerGetter struct{}

func (g *SBMServerGetter) Get(ctx context.Context, client *serverscom.Client, id string) (interface{}, error) {
	return client.Hosts.GetSBMServer(ctx, id)
}

func newGetCmd(cmdContext *base.CmdContext, hostType *HostType) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: fmt.Sprintf("Get a %s", hostType.entityName),
		Long:  fmt.Sprintf("Get a %s by id", hostType.entityName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := hostType.getter.Get(ctx, scClient, id)
			if err != nil {
				return err
			}

			if entity != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(entity)
			}
			return nil
		},
	}
	return cmd
}
