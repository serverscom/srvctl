package ssh

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/serverscom/srvctl/internal/config"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	manager, err := config.NewManager()
	if err != nil {
		log.Fatal(err)
	}

	factory := func(verbose bool) serverscom.Collection[serverscom.SSHKey] {
		c := client.NewClient(manager.GetToken(), manager.GetEndpoint()).SetVerbose(verbose)
		return c.GetScClient().SSHKeys.Collection()
	}

	return base.NewListCmd("SSH Keys", factory, manager)
}
