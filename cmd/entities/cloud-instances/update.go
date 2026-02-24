package cloudinstances

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpdateCmd(cmdContext *base.CmdContext) *cobra.Command {
	var name string
	var backupCopies int
	var ipv6Enabled bool
	var gpnEnabled bool
	var labels []string

	cmd := &cobra.Command{
		Use:   "update <instance-id>",
		Short: "Update cloud instance by instance ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			labelsMap, err := base.ParseLabels(labels)
			if err != nil {
				return err
			}

			input := serverscom.CloudComputingInstanceUpdateInput{}
			if cmd.Flags().Changed("name") {
				input.Name = &name
			}
			if cmd.Flags().Changed("backup-copies") {
				input.BackupCopies = &backupCopies
			}
			if cmd.Flags().Changed("ipv6-enabled") {
				input.IPv6Enabled = &ipv6Enabled
			}
			if cmd.Flags().Changed("gpn-enabled") {
				input.GPNEnabled = &gpnEnabled
			}
			if cmd.Flags().Changed("label") {
				input.Labels = labelsMap
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			instanceID := args[0]
			out, err := scClient.CloudComputingInstances.Update(ctx, instanceID, input)
			if err != nil {
				return err
			}

			if out != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(out)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the cloud instance")
	cmd.Flags().IntVar(&backupCopies, "backup-copies", 0, "Number of backup copies")
	cmd.Flags().BoolVar(&ipv6Enabled, "ipv6-enabled", false, "Enable IPv6")
	cmd.Flags().BoolVar(&gpnEnabled, "gpn-enabled", false, "Enable GPN")
	cmd.Flags().StringArrayVar(&labels, "label", []string{}, "Labels in key=value format")

	return cmd
}
