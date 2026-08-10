package cmd

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/cmd/config"
	"github.com/serverscom/srvctl/cmd/context"
	"github.com/serverscom/srvctl/cmd/entities/account"
	cloudbackups "github.com/serverscom/srvctl/cmd/entities/cloud-backups"
	cloudinstances "github.com/serverscom/srvctl/cmd/entities/cloud-instances"
	cloudregions "github.com/serverscom/srvctl/cmd/entities/cloud-regions"
	cloudvolumes "github.com/serverscom/srvctl/cmd/entities/cloud-volumes"
	"github.com/serverscom/srvctl/cmd/entities/drivemodels"
	"github.com/serverscom/srvctl/cmd/entities/hosts"
	"github.com/serverscom/srvctl/cmd/entities/invoices"
	"github.com/serverscom/srvctl/cmd/entities/k8s"
	l2segments "github.com/serverscom/srvctl/cmd/entities/l2_segments"
	loadbalancerclusters "github.com/serverscom/srvctl/cmd/entities/load_balancer_clusters"
	loadbalancers "github.com/serverscom/srvctl/cmd/entities/load_balancers"
	"github.com/serverscom/srvctl/cmd/entities/locations"
	networkpools "github.com/serverscom/srvctl/cmd/entities/network-pools"
	"github.com/serverscom/srvctl/cmd/entities/racks"
	rbsvolumes "github.com/serverscom/srvctl/cmd/entities/rbs_volumes"
	sbmmodels "github.com/serverscom/srvctl/cmd/entities/sbm_models"
	sbmosoptions "github.com/serverscom/srvctl/cmd/entities/sbm_os_options"
	serverosoptions "github.com/serverscom/srvctl/cmd/entities/server_os_options"
	serverramoptions "github.com/serverscom/srvctl/cmd/entities/server_ram_options"
	"github.com/serverscom/srvctl/cmd/entities/servermodels"
	sshkeys "github.com/serverscom/srvctl/cmd/entities/ssh-keys"
	"github.com/serverscom/srvctl/cmd/entities/ssl"
	"github.com/serverscom/srvctl/cmd/entities/uplinkbandwidths"
	"github.com/serverscom/srvctl/cmd/entities/uplinkmodels"
	"github.com/serverscom/srvctl/cmd/login"
	"github.com/serverscom/srvctl/internal/client"
	"github.com/spf13/cobra"
)

const (
	groupResources = "resources"
	groupConfig    = "config"
	groupOther     = "other"
)

func addGroupedCommands(parent *cobra.Command, groupID string, cmds ...*cobra.Command) {
	for _, c := range cmds {
		c.GroupID = groupID
		parent.AddCommand(c)
	}
}

func NewRootCmd(version string) *cobra.Command {
	cobra.EnableTraverseRunHooks = true

	cmdContext := &base.CmdContext{}

	cmd := &cobra.Command{
		Use:               "srvctl [command] [flags]",
		Short:             "CLI tool for servers.com API",
		Long:              `A command line interface for managing servers.com resources`,
		Version:           version,
		PersistentPreRunE: base.InitCmdContext(cmdContext),
		SilenceUsage:      true,
	}
	// Global flags
	base.AddGlobalFlags(cmd)

	cmd.AddGroup(
		&cobra.Group{ID: groupResources, Title: "Resource Commands:"},
		&cobra.Group{ID: groupConfig, Title: "Configuration Commands:"},
		&cobra.Group{ID: groupOther, Title: "Other Commands:"},
	)

	clientFactory := &client.DefaultClientFactory{}

	// Configuration/CLI commands
	addGroupedCommands(cmd, groupConfig,
		login.NewCmd(cmdContext, clientFactory),
		context.NewCmd(cmdContext),
		config.NewCmd(cmdContext),
	)

	// Resource commands
	addGroupedCommands(cmd, groupResources,
		sshkeys.NewCmd(cmdContext),
		hosts.NewCmd(cmdContext),
		ssl.NewCmd(cmdContext),
		loadbalancerclusters.NewCmd(cmdContext),
		loadbalancers.NewCmd(cmdContext),
		racks.NewCmd(cmdContext),
		invoices.NewCmd(cmdContext),
		account.NewCmd(cmdContext),
		locations.NewCmd(cmdContext),
		k8s.NewCmd(cmdContext),
		uplinkmodels.NewCmd(cmdContext),
		uplinkbandwidths.NewCmd(cmdContext),
		servermodels.NewCmd(cmdContext),
		drivemodels.NewCmd(cmdContext),
		serverosoptions.NewCmd(cmdContext),
		serverramoptions.NewCmd(cmdContext),
		sbmosoptions.NewCmd(cmdContext),
		sbmmodels.NewCmd(cmdContext),
		l2segments.NewCmd(cmdContext),
		networkpools.NewCmd(cmdContext),
		cloudinstances.NewCmd(cmdContext),
		cloudregions.NewCmd(cmdContext),
		cloudvolumes.NewCmd(cmdContext),
		cloudbackups.NewCmd(cmdContext),
		rbsvolumes.NewCmd(cmdContext),
	)

	cmd.SetHelpCommandGroupID(groupOther)
	cmd.SetCompletionCommandGroupID(groupOther)

	return cmd
}
