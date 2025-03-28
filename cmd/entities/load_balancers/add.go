package loadbalancers

import (
	"context"
	"fmt"
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type LBCreator interface {
	Create(ctx context.Context, client *serverscom.Client, input any) (any, error)
	NewCreateInput() any
}

type LBL4CreateMgr struct{}

func (c *LBL4CreateMgr) Create(ctx context.Context, client *serverscom.Client, input any) (any, error) {
	lbInput, ok := input.(*serverscom.L4LoadBalancerCreateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for L4 LB")
	}
	return client.LoadBalancers.CreateL4LoadBalancer(ctx, *lbInput)
}

func (c *LBL4CreateMgr) NewCreateInput() any {
	return &serverscom.L4LoadBalancerCreateInput{}
}

type LBL7CreateMgr struct{}

func (c *LBL7CreateMgr) Create(ctx context.Context, client *serverscom.Client, input any) (any, error) {
	lbInput, ok := input.(*serverscom.L7LoadBalancerCreateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for L7 LB")
	}
	return client.LoadBalancers.CreateL7LoadBalancer(ctx, *lbInput)
}

func (c *LBL7CreateMgr) NewCreateInput() any {
	return &serverscom.L7LoadBalancerCreateInput{}
}

func newAddCmd(cmdContext *base.CmdContext, lbType *LBTypeCmd) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: fmt.Sprintf("Create %s", lbType.entityName),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := lbType.managers.createMgr.NewCreateInput()

			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			lb, err := lbType.managers.createMgr.Create(ctx, scClient, input)
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
