package sshkeys

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.SSHKey] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.SSHKeys.Collection()
	}

	opts := &base.BaseLabelsListOptions[serverscom.SSHKey]{}

	return base.NewListCmd("list", "SSH Keys", factory, cmdContext, opts)
}
