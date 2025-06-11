package k8s

import (
	"log"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/base"
	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/spf13/cobra"
)

func NewCmd(cmdContext *base.CmdContext) *cobra.Command {
	entitiesMap, err := getK8sEntities()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "Manage kubernetes clusters",
		PersistentPreRunE: base.CombinePreRunE(
			base.CheckFormatterFlags(cmdContext, entitiesMap),
			base.CheckEmptyContexts(cmdContext),
		),
		Args: base.NoArgs,
		Run:  base.UsageRun,
	}

	cmd.AddCommand(
		newListCmd(cmdContext),
		newListNodesCmd(cmdContext),
		newGetCmd(cmdContext),
		newGetNodeCmd(cmdContext),
		newUpdateCmd(cmdContext),
	)

	base.AddFormatFlags(cmd)

	return cmd
}

func getK8sEntities() (map[string]entities.EntityInterface, error) {
	result := make(map[string]entities.EntityInterface)

	k8sEntity, err := entities.Registry.GetEntityFromValue(serverscom.KubernetesCluster{})
	if err != nil {
		log.Fatal(err)
	}
	result["k8s"] = k8sEntity

	k8sNodeEntity, err := entities.Registry.GetEntityFromValue(serverscom.KubernetesClusterNode{})
	if err != nil {
		log.Fatal(err)
	}
	result["list-nodes"] = k8sNodeEntity
	result["get-node"] = k8sNodeEntity

	return result, nil
}
