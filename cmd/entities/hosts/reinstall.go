package hosts

import (
	"context"
	"fmt"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type AddedFlags struct {
	Skeleton  bool
	InputPath string
}

type HostReinstaller interface {
	Reinstall(ctx context.Context, client *serverscom.Client, id string, input any) (any, error)
	NewReinstallInput() any
}

type DSReinstallMgr struct{}

func (r *DSReinstallMgr) Reinstall(ctx context.Context, client *serverscom.Client, id string, input any) (any, error) {
	dsInput, ok := input.(*serverscom.OperatingSystemReinstallInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for dedicated server")
	}
	return client.Hosts.ReinstallOperatingSystemForDedicatedServer(ctx, id, *dsInput)
}

func (c *DSReinstallMgr) NewReinstallInput() any {
	return &serverscom.OperatingSystemReinstallInput{}
}

type SBMReinstallMgr struct{}

func (r *SBMReinstallMgr) Reinstall(ctx context.Context, client *serverscom.Client, id string, input any) (any, error) {
	sbmInput, ok := input.(*serverscom.SBMOperatingSystemReinstallInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type for SBM server")
	}
	return client.Hosts.ReinstallOperatingSystemForSBMServer(ctx, id, *sbmInput)
}

func (c *SBMReinstallMgr) NewReinstallInput() any {
	return &serverscom.SBMOperatingSystemReinstallInput{}
}

func newReinstallCmd(cmdContext *base.CmdContext, hostType *HostTypeCmd) *cobra.Command {
	flags := &AddedFlags{}

	cmd := &cobra.Command{
		Use:   "reinstall <id>",
		Short: fmt.Sprintf("Reinstall OS for a  %s", hostType.entityName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()
			base.SetupProxy(cmd, manager)

			input := hostType.managers.reinstallMgr.NewReinstallInput()

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), &input); err != nil {
					return err
				}
			} else if flags.Skeleton {
				return formatter.FormatSkeleton("hosts/reinstall.json")
			} else {
				required := []string{"input"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			id := args[0]

			server, err := hostType.managers.reinstallMgr.Reinstall(ctx, scClient, id, input)
			if err != nil {
				return err
			}

			if server != nil {
				return formatter.Format(server)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	return cmd
}
