package sshkeys

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool) serverscom.Collection[serverscom.SSHKey] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.SSHKeys.Collection()
	}

	opts := &base.BaseListOptions[serverscom.SSHKey]{}

	return base.NewListCmd("SSH Keys", factory, cmdContext, opts)
}
