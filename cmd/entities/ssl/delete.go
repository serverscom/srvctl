package ssl

import (
	"context"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type SSLDeleter interface {
	Delete(ctx context.Context, client *serverscom.Client, id string) error
}

type SSLCustomDeleteMgr struct{}

func (d *SSLCustomDeleteMgr) Delete(ctx context.Context, client *serverscom.Client, id string) error {
	return client.SSLCertificates.DeleteCustom(ctx, id)
}

type SSLLeDeleteMgr struct{}

func (d *SSLLeDeleteMgr) Delete(ctx context.Context, client *serverscom.Client, id string) error {
	return client.SSLCertificates.DeleteLE(ctx, id)
}

func newDeleteCmd(cmdContext *base.CmdContext, sslType *SSLTypeCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an SSL certificate",
		Long:  "Delete an SSL certificate by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			err := sslType.managers.deleteMgr.Delete(ctx, scClient, id)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
