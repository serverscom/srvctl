package hosts

import (
	"context"
	"fmt"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type HostGetter interface {
	Get(ctx context.Context, client *serverscom.Client, id string) (any, error)
}

type DSGetMgr struct{}

func (g *DSGetMgr) Get(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.Hosts.GetDedicatedServer(ctx, id)
}

type KBMGetMgr struct{}

func (g *KBMGetMgr) Get(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.Hosts.GetKubernetesBaremetalNode(ctx, id)
}

type SBMGetMgr struct{}

func (g *SBMGetMgr) Get(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.Hosts.GetSBMServer(ctx, id)
}

func newGetCmd(cmdContext *base.CmdContext, hostType *HostTypeCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: fmt.Sprintf("Get a %s", hostType.entityName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			entity, err := hostType.managers.getMgr.Get(ctx, scClient, id)
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
