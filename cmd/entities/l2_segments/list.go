package l2segments

import (
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

func newListCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.L2Segment] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.L2Segments.Collection()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.L2Segment]{},
		&base.LabelSelectorOption[serverscom.L2Segment]{},
	)

	return base.NewListCmd("list", "l2", factory, cmdContext, opts...)
}

func newListGroupsCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.L2LocationGroup] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.L2Segments.LocationGroups()
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.L2LocationGroup]{},
		&base.L2SegmentGroupTypeOption[serverscom.L2LocationGroup]{},
	)

	return base.NewListCmd("list-groups", "l2 location groups", factory, cmdContext, opts...)
}

func newListMembersCmd(cmdContext *base.CmdContext) *cobra.Command {
	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.L2Member] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.L2Segments.Members(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.L2Member]{},
		&base.L2SegmentGroupTypeOption[serverscom.L2Member]{},
	)

	return base.NewListCmd("list-members", "l2 members", factory, cmdContext, opts...)
}

func newListNetworksCmd(cmdContext *base.CmdContext) *cobra.Command {

	factory := func(verbose bool, args ...string) serverscom.Collection[serverscom.Network] {
		scClient := cmdContext.GetClient().SetVerbose(verbose).GetScClient()
		return scClient.L2Segments.Networks(args[0])
	}

	opts := base.NewListOptions(
		&base.BaseListOptions[serverscom.Network]{},
		&base.L2SegmentGroupTypeOption[serverscom.Network]{},
	)

	return base.NewListCmd("list-networks", "l2 networks", factory, cmdContext, opts...)
}
