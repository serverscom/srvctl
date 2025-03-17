package ssl

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.SSLCertificate] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.SSLCertificates.Collection()
	}

	opts := &base.BaseLabelsListOptions[serverscom.SSLCertificate]{}

	return base.NewListCmd("list", "SSL Certificates", factory, cmdContext, opts)
}
