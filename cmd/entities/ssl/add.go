package ssl

import (
	"context"
	"fmt"
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type SSLCreator interface {
	Create(ctx context.Context, client *serverscom.Client, input any) (any, error)
	NewCreateInput() any
}

type SSLCustomCreateMgr struct{}

func (c *SSLCustomCreateMgr) Create(ctx context.Context, client *serverscom.Client, input any) (any, error) {
	sslInput, ok := input.(*serverscom.SSLCertificateCreateCustomInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for custom SSL")
	}
	return client.SSLCertificates.CreateCustom(ctx, *sslInput)
}

func (c *SSLCustomCreateMgr) NewCreateInput() any {
	return &serverscom.SSLCertificateCreateCustomInput{}
}

func newAddCmd(cmdContext *base.CmdContext, sslType *SSLTypeCmd) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: fmt.Sprintf("Create a %s", sslType.entityName),
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := sslType.managers.createMgr.NewCreateInput()

			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			server, err := sslType.managers.createMgr.Create(ctx, scClient, input)
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
