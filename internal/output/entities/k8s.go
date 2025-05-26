package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	KubernetesClusterType                  = reflect.TypeOf(serverscom.KubernetesCluster{})
	KubernetesClusterNodeType              = reflect.TypeOf(serverscom.KubernetesClusterNode{})
	KubernetesClusterListDefaultFields     = []string{"ID", "Name", "Status", "LocationID"}
	KubernetesClusterNodeListDefaultFields = []string{"ID", "Number", "Hostname", "Type", "Role", "Status", "PrivateIPv4Address", "PublicIPv4Address"}
)

func RegisterKubernetesClusterDefinition() {
	k8sEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": KubernetesClusterListDefaultFields,
		},
		eType: KubernetesClusterType,
	}
	if err := Registry.Register(k8sEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterKubernetesClusterNodeDefinition() {
	k8sNodeEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Number", Name: "Number", Path: "Number", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Hostname", Name: "Hostname", Path: "Hostname", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Role", Name: "Role", Path: "Role", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RefID", Name: "RefID", Path: "RefID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ClusterID", Name: "ClusterID", Path: "ClusterID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list-nodes": KubernetesClusterNodeListDefaultFields,
		},
		eType: KubernetesClusterNodeType,
	}
	if err := Registry.Register(k8sNodeEntity); err != nil {
		log.Fatal(err)
	}
}
