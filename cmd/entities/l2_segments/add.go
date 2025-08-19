package l2segments

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "add --input <path>",
		Short: "Add a new L2 segment",
		Long:  "Add a new L2 segment",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.L2SegmentCreateInput{}
			if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
				return err
			}

			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()
			l2Segment, err := scClient.L2Segments.Create(ctx, *input)
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
