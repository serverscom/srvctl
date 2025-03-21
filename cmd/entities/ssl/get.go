package ssl

import (
	"context"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type SSLGetter interface {
	Get(ctx context.Context, client *serverscom.Client, id string) (any, error)
}

type SSLCustomGetMgr struct{}

func (g *SSLCustomGetMgr) Get(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.SSLCertificates.GetCustom(ctx, id)
}

type SSLLeGetMgr struct{}

func (g *SSLLeGetMgr) Get(ctx context.Context, client *serverscom.Client, id string) (any, error) {
	return client.SSLCertificates.GetLE(ctx, id)
}

func newGetCmd(cmdContext *base.CmdContext, sslType *SSLTypeCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an ssl certificate",
		Long:  "Get an ssl certificate by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			sslCert, err := sslType.managers.getMgr.Get(ctx, scClient, id)
			if err != nil {
				return err
			}

			if sslCert != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(sslCert)
			}
			return nil
		},
	}
	return cmd
}
