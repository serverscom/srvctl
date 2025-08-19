package l2segments

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newUpdateL2Cmd(cmdContext *base.CmdContext) *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "update <l2_segment_id>",
		Short: "Update an L2 segment",
		Long:  "Update an L2 segment by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.L2SegmentUpdateInput{}
			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			l2SegmentId := args[0]
			l2Segment, err := scClient.L2Segments.Update(ctx, l2SegmentId, *input)
			if err != nil {
				return err
			}

			if l2Segment != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(l2Segment)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "input", "i", "", "path to input file or '-' to read from stdin")
	if err := cmd.MarkFlagRequired("input"); err != nil {
		log.Fatal(err)
	}

	return cmd
}

func newUpdateL2NetworksCmd(cmdContext *base.CmdContext) *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "update-networks <l2_segment_id>",
		Short: "Update an L2 segment networks",
		Long:  "Update an L2 segment networks by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.L2SegmentChangeNetworksInput{}
			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

			l2SegmentId := args[0]
			l2Segment, err := scClient.L2Segments.ChangeNetworks(ctx, l2SegmentId, *input)
			if err != nil {
				return err
			}

			if l2Segment != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(l2Segment)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "input", "i", "", "path to input file or '-' to read from stdin")
	if err := cmd.MarkFlagRequired("input"); err != nil {
		log.Fatal(err)
	}

	return cmd
}
