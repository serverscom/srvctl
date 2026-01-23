package l2segments

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type UpdateFlags struct {
	Skeleton  bool
	InputPath string
}

func newUpdateL2Cmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &UpdateFlags{}

	cmd := &cobra.Command{
		Use:   "update <l2_segment_id>",
		Short: "Update an L2 segment",
		Long:  "Update an L2 segment by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.L2SegmentUpdateInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else if flags.Skeleton {
				return formatter.FormatSkeleton("l2-segments/update.json")
			} else {
				required := []string{"input"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			l2SegmentId := args[0]
			l2Segment, err := scClient.L2Segments.Update(ctx, l2SegmentId, *input)
			if err != nil {
				return err
			}

			if l2Segment != nil {
				return formatter.Format(l2Segment)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	return cmd
}

func newUpdateL2NetworksCmd(cmdContext *base.CmdContext) *cobra.Command {
	flags := &UpdateFlags{}

	cmd := &cobra.Command{
		Use:   "update-networks <l2_segment_id>",
		Short: "Update an L2 segment networks",
		Long:  "Update an L2 segment networks by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			formatter := cmdContext.GetOrCreateFormatter(cmd)
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.L2SegmentChangeNetworksInput{}

			if flags.InputPath != "" {
				if err := base.ReadInputJSON(flags.InputPath, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else if flags.Skeleton {
				return formatter.FormatSkeleton("l2-segments/update_networks.json")
			} else {
				required := []string{"input"}
				if err := base.ValidateFlags(cmd, required); err != nil {
					return err
				}
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			l2SegmentId := args[0]
			l2Segment, err := scClient.L2Segments.ChangeNetworks(ctx, l2SegmentId, *input)
			if err != nil {
				return err
			}

			if l2Segment != nil {
				return formatter.Format(l2Segment)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.InputPath, "input", "i", "", "path to input file or '-' to read from stdin")
	cmd.Flags().BoolVarP(&flags.Skeleton, "skeleton", "s", false, "JSON object with structure that is required to be passed")

	return cmd
}
