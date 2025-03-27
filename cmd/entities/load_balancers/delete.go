package loadbalancers

import (
	"context"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type LBDeleter interface {
	Delete(ctx context.Context, client *serverscom.Client, id string) error
}

type LBL4DeleteMgr struct{}

func (d *LBL4DeleteMgr) Delete(ctx context.Context, client *serverscom.Client, id string) error {
	return client.LoadBalancers.DeleteL4LoadBalancer(ctx, id)
}

type LBL7DeleteMgr struct{}

func (d *LBL7DeleteMgr) Delete(ctx context.Context, client *serverscom.Client, id string) error {
	return client.LoadBalancers.DeleteL7LoadBalancer(ctx, id)
}

func newDeleteCmd(cmdContext *base.CmdContext, lbType *LBTypeCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete load balancer",
		Long:  "Delete load balancer by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			err := lbType.managers.deleteMgr.Delete(ctx, scClient, id)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
