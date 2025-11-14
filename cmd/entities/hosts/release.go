package hosts

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newDSAbortReleaseCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "abort-release <id>",
		Short: "Abort release for a dedicated server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			server, err := scClient.Hosts.AbortReleaseForDedicatedServer(ctx, id)
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
	return cmd
}

func newDSScheduleReleaseCmd(cmdContext *base.CmdContext) *cobra.Command {
	var releaseAfter string

	cmd := &cobra.Command{
		Use:   "schedule-release <id>",
		Short: "Schedule release for a dedicated server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			input := serverscom.ScheduleReleaseInput{ReleaseAfter: releaseAfter}
			server, err := scClient.Hosts.ScheduleReleaseForDedicatedServer(ctx, id, input)
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

	cmd.Flags().StringVar(&releaseAfter, "release-after", "", "UTC datetime string in format: YYYY-MM-DDTHH:MM:SS+HH:MM")

	return cmd
}

func newSBMReleaseCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <id>",
		Short: "Release an SBM server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)
			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			id := args[0]
			server, err := scClient.Hosts.ReleaseSBMServer(ctx, id)
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

	return cmd
}
