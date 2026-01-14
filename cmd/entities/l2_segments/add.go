package l2segments

import (
	"fmt"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
	"strings"
)

var (
	path            string
	name            string
	l2Type          string
	locationGroupID int64
	members         []string
	labels          []string
)

func newAddCmd(cmdContext *base.CmdContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new L2 segment",
		Long:  "Add a new L2 segment",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()

			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			input := &serverscom.L2SegmentCreateInput{}

			if cmd.Flags().Changed("input") {
				if err := base.ReadInputJSON(path, cmd.InOrStdin(), input); err != nil {
					return err
				}
			} else {
				if err := validateFlags(cmd); err != nil {
					return err
				}
				if err := fillInput(cmd, input); err != nil {
					return err
				}
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
	cmd.Flags().StringVarP(&name, "name", "n", "", "A name of a L2 segment")
	cmd.Flags().StringVarP(&l2Type, "type", "", "", "A type of a L2 segment")
	cmd.Flags().Int64VarP(&locationGroupID, "location-group-id", "", 0, "A private-key of a L2 segment")
	cmd.Flags().StringArrayVarP(&members, "member", "m", []string{}, "L2 segment member: id=<string>,mode=<native|trunk>")
	cmd.Flags().StringArrayVarP(&labels, "label", "l", []string{}, "string in key=value format")

	return cmd
}

func validateFlags(cmd *cobra.Command) error {
	required := []string{"name", "type", "location-group-id", "member"}
	var missing []string

	for _, flag := range required {
		if !cmd.Flags().Changed(flag) {
			missing = append(missing, "--"+flag)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"use --input or provide all required flags (missing: %s)",
			strings.Join(missing, ", "),
		)
	}

	return nil
}

func fillInput(cmd *cobra.Command, input *serverscom.L2SegmentCreateInput) error {
	if cmd.Flags().Changed("name") {
		input.Name = &name
	}
	if cmd.Flags().Changed("type") {
		input.Type = l2Type
	}
	if cmd.Flags().Changed("location-group-id") {
		input.LocationGroupID = locationGroupID
	}
	if cmd.Flags().Changed("member") {
		membersMap, err := parseMembers(members)
		if err != nil {
			return err
		}

		input.Members = membersMap
	}
	if cmd.Flags().Changed("label") {
		labelsMap, err := base.ParseLabels(labels)
		if err != nil {
			return err
		}

		input.Labels = labelsMap
	}

	return nil
}

func parseMembers(members []string) ([]serverscom.L2SegmentMemberInput, error) {
	var res []serverscom.L2SegmentMemberInput

	for _, member := range members {
		m := serverscom.L2SegmentMemberInput{}
		parts := strings.Split(member, ",")

		for _, p := range parts {
			props := strings.SplitN(p, "=", 2)
			if len(props) != 2 {
				return nil, fmt.Errorf("invalid member format: %s", p)
			}

			switch props[0] {
			case "id":
				m.ID = props[1]
			case "mode":
				m.Mode = props[1]
			default:
				return nil, fmt.Errorf("unknown member field: %s", props[0])
			}
		}

		if m.ID == "" || m.Mode == "" {
			return nil, fmt.Errorf("member must include id and mode: %s", member)
		}

		res = append(res, m)
	}

	return res, nil
}
