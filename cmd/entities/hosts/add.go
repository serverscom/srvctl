package hosts

import (
	"context"
	"fmt"
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type HostCreator interface {
	Create(ctx context.Context, client *serverscom.Client, input any) (any, error)
	NewCreateInput() any
}

type DSCreateMgr struct{}

func (c *DSCreateMgr) Create(ctx context.Context, client *serverscom.Client, input any) (any, error) {
	dsInput, ok := input.(*serverscom.DedicatedServerCreateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for dedicated server")
	}
	return client.Hosts.CreateDedicatedServers(ctx, *dsInput)
}

func (c *DSCreateMgr) NewCreateInput() any {
	return &serverscom.DedicatedServerCreateInput{}
}

type SBMCreateMgr struct{}

func (c *SBMCreateMgr) Create(ctx context.Context, client *serverscom.Client, input any) (any, error) {
	sbmInput, ok := input.(*serverscom.SBMServerCreateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for SBM server")
	}
	return client.Hosts.CreateSBMServers(ctx, *sbmInput)
}

func (c *SBMCreateMgr) NewCreateInput() any {
	return &serverscom.SBMServerCreateInput{}
}

func newAddCmd(cmdContext *base.CmdContext, hostType *HostTypeCmd) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: fmt.Sprintf("Create a %s", hostType.entityName),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := hostType.managers.createMgr.NewCreateInput()

			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			server, err := hostType.managers.createMgr.Create(ctx, scClient, input)
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

	cmd.Flags().StringVarP(&path, "input", "i", "", "path to input file or '-' to read from stdin")
	if err := cmd.MarkFlagRequired("input"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
