package loadbalancers

import (
	"context"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type LBGetter interface {
	Get(ctx context.Context, client *serverscom.Client, id string) (any, error)
}

type LBL4GetMgr struct{}

func (g *LBL4GetMgr) Get(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.LoadBalancers.GetL4LoadBalancer(ctx, id)
}

type LBL7GetMgr struct{}

func (g *LBL7GetMgr) Get(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.LoadBalancers.GetL7LoadBalancer(ctx, id)
}

func newGetCmd(cmdContext *base.CmdContext, lbType *LBTypeCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get load balancer",
		Long:  "Get load balancer by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			lb, err := lbType.managers.getMgr.Get(ctx, scClient, id)
			if err != nil {
				return err
			}

			if lb != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(lb)
			}
			return nil
		},
	}
	return cmd
}
