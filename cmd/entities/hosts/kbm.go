package hosts

import (
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

// NewKBMCmd builds the top-level "kbm" command.
func NewKBMCmd(cmdContext *base.CmdContext) *cobra.Command {
	return newHostTypeCmd(cmdContext, HostTypeCmd{
		use:        "kbm",
		shortDesc:  "Manage kubernetes baremetal nodes",
		entityName: "Kubernetes baremetal nodes",
		typeFlag:   "kubernetes_baremetal_node",
		managers: HostManagers{
			getMgr:   &KBMGetMgr{},
			powerMgr: &KBMPowerMgr{},
		},
		extraCmds: []func(*base.CmdContext) *cobra.Command{
			newUpdateKBMCmd,
			newListKBMCmd,
			newListKBMNetworksCmd,
			newListKBMDriveSlotsCmd,
		},
	})
}
