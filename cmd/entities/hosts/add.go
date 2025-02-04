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
	Create(ctx context.Context, client *serverscom.Client, input interface{}) (interface{}, error)
	NewCreateInput() interface{}
}

type DedicatedServerCreator struct{}

func (c *DedicatedServerCreator) Create(ctx context.Context, client *serverscom.Client, input interface{}) (interface{}, error) {
	dsInput, ok := input.(*serverscom.DedicatedServerCreateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for dedicated server")
	}
	return client.Hosts.CreateDedicatedServers(ctx, *dsInput)
}

func (c *DedicatedServerCreator) NewCreateInput() interface{} {
	return &serverscom.DedicatedServerCreateInput{}
}

type SBMServerCreator struct{}

func (c *SBMServerCreator) Create(ctx context.Context, client *serverscom.Client, input interface{}) (interface{}, error) {
	sbmInput, ok := input.(*serverscom.SBMServerCreateInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for SBM server")
	}
	return client.Hosts.CreateSBMServers(ctx, *sbmInput)
}

func (c *SBMServerCreator) NewCreateInput() interface{} {
	return &serverscom.SBMServerCreateInput{}
}

func newAddCmd(cmdContext *base.CmdContext, hostType *HostType) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: fmt.Sprintf("Create a new %s", hostType.entityName),
		Long:  fmt.Sprintf("Create a new %s", hostType.entityName),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := hostType.creator.NewCreateInput()

			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			server, err := hostType.creator.Create(ctx, scClient, input)
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
