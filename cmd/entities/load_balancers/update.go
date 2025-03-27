package loadbalancers

import (
	"context"
	"fmt"
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type LBUpdater interface {
	Update(ctx context.Context, client *serverscom.Client, id string, input any) (any, error)
	NewUpdateInput() any
}

type LBL4UpdateMgr struct{}

func (c *LBL4UpdateMgr) Update(ctx context.Context, client *serverscom.Client, id string, input any) (any, error) {
	lbInput, ok := input.(*serverscom.L4LoadBalancerUpdateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for L4 LB")
	}
	return client.LoadBalancers.UpdateL4LoadBalancer(ctx, id, *lbInput)
}

func (c *LBL4UpdateMgr) NewUpdateInput() any {
	return &serverscom.L4LoadBalancerUpdateInput{}
}

type LBL7UpdateMgr struct{}

func (c *LBL7UpdateMgr) Update(ctx context.Context, client *serverscom.Client, id string, input any) (any, error) {
	lbInput, ok := input.(*serverscom.L7LoadBalancerUpdateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for L7 LB")
	}
	return client.LoadBalancers.UpdateL7LoadBalancer(ctx, id, *lbInput)
}

func (c *LBL7UpdateMgr) NewUpdateInput() any {
	return &serverscom.L7LoadBalancerUpdateInput{}
}

func newUpdateCmd(cmdContext *base.CmdContext, lbType *LBTypeCmd) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "update --input <path>",
		Short: fmt.Sprintf("Update %s", lbType.entityName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := lbType.managers.updateMgr.NewUpdateInput()

			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			lb, err := lbType.managers.updateMgr.Update(ctx, scClient, id, input)
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

	cmd.Flags().StringVarP(&path, "input", "i", "", "path to input file or '-' to read from stdin")
	if err := cmd.MarkFlagRequired("input"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
