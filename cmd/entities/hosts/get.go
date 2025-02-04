package hosts

import (
	"context"
	"fmt"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/spf13/cobra"
)

type HostGetter interface {
	Get(ctx context.Context, id string) (interface{}, error)
}

type DedicatedServerGetter struct {
	client *serverscom.Client
}

func (g *DedicatedServerGetter) Get(ctx context.Context, id string) (interface{}, error) {
	return g.client.Hosts.GetDedicatedServer(ctx, id)
}

type KubernetesBaremetalNodeGetter struct {
	client *serverscom.Client
}

func (g *KubernetesBaremetalNodeGetter) Get(ctx context.Context, id string) (interface{}, error) {
	return g.client.Hosts.GetKubernetesBaremetalNode(ctx, id)
}

type SBMServerGetter struct {
	client *serverscom.Client
}

func (g *SBMServerGetter) Get(ctx context.Context, id string) (interface{}, error) {
	return g.client.Hosts.GetSBMServer(ctx, id)
}

func newGetCmd(cmdContext *base.CmdContext, hostType *HostType) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: fmt.Sprintf("Get a %s", hostType.entityName),
		Long:  fmt.Sprintf("Get a %s by id", hostType.entityName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cmdContext.GetManager()
			ctx, cancel := base.SetupContext(cmd, manager)
			defer cancel()

			base.SetupProxy(cmd, manager)

			id := args[0]
			entity, err := hostType.getter.Get(ctx, id)
			if err != nil {
				return err
			}

			if entity != nil {
				formatter := cmdContext.GetOrCreateFormatter(cmd)
				return formatter.Format(entity)
			}
			return nil
		},
	}
	return cmd
}

// func newGetDsCmd(cmdContext *base.CmdContext) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "get <id>",
// 		Short: "Get a dedicated server",
// 		Long:  "Get a dedicated server by id",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			manager := cmdContext.GetManager()

// 			ctx, cancel := base.SetupContext(cmd, manager)
// 			defer cancel()

// 			base.SetupProxy(cmd, manager)

// 			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

// 			id := args[0]
// 			ds, err := scClient.Hosts.GetDedicatedServer(ctx, id)
// 			if err != nil {
// 				return err
// 			}

// 			if ds != nil {
// 				formatter := cmdContext.GetOrCreateFormatter(cmd)
// 				return formatter.Format(ds)
// 			}
// 			return nil
// 		},
// 	}

// 	return cmd
// }

// func newGetKbmCmd(cmdContext *base.CmdContext) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "get <id>",
// 		Short: "Get a kubernetes baremetal node",
// 		Long:  "Get a kubernetes baremetal node by id",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			manager := cmdContext.GetManager()

// 			ctx, cancel := base.SetupContext(cmd, manager)
// 			defer cancel()

// 			base.SetupProxy(cmd, manager)

// 			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

// 			id := args[0]
// 			kbm, err := scClient.Hosts.GetKubernetesBaremetalNode(ctx, id)
// 			if err != nil {
// 				return err
// 			}

// 			if kbm != nil {
// 				formatter := cmdContext.GetOrCreateFormatter(cmd)
// 				return formatter.Format(kbm)
// 			}
// 			return nil
// 		},
// 	}

// 	return cmd
// }

// func newGetSbmCmd(cmdContext *base.CmdContext) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "get <id>",
// 		Short: "Get a scalable baremetal server",
// 		Long:  "Get a scalable baremetal server by id",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			manager := cmdContext.GetManager()

// 			ctx, cancel := base.SetupContext(cmd, manager)
// 			defer cancel()

// 			base.SetupProxy(cmd, manager)

// 			scClient := cmdContext.GetClient().SetVerbose(manager.GetVerbose(cmd)).GetScClient()

// 			id := args[0]
// 			sbm, err := scClient.Hosts.GetSBMServer(ctx, id)
// 			if err != nil {
// 				return err
// 			}

// 			if sbm != nil {
// 				formatter := cmdContext.GetOrCreateFormatter(cmd)
// 				return formatter.Format(sbm)
// 			}
// 			return nil
// 		},
// 	}

// 	return cmd
// }
