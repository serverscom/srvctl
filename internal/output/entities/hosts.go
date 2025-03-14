package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	HostType                    = reflect.TypeOf(serverscom.Host{})
	DedicatedServerType         = reflect.TypeOf(serverscom.DedicatedServer{})
	KubernetesBaremetalNodeType = reflect.TypeOf(serverscom.KubernetesBaremetalNode{})
	SBMServerType               = reflect.TypeOf(serverscom.SBMServer{})
	HostListDefaultFields       = []string{"ID", "Type", "Title", "Status"}
	CmdDefaultFields            = map[string][]string{
		"list": HostListDefaultFields,
	}
)

func getConfigurationDetailsField() Field {
	f := Field{
		ID:                  "ConfigurationDetails",
		Name:                "ConfigurationDetails",
		Path:                "ConfigurationDetails",
		PageViewHandlerFunc: structPVHandler,
	}
	childs := []Field{
		{ID: "RAMSize", Name: "RAMSize", Path: "RAMSize", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "ServerModelID", Name: "ServerModelID", Path: "ServerModelID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "ServerModelName", Name: "ServerModelName", Path: "ServerModelName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "PublicUplinkID", Name: "PublicUplinkID", Path: "PublicUplinkID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "PublicUplinkName", Name: "PublicUplinkName", Path: "PublicUplinkName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "PrivateUplinkID", Name: "PrivateUplinkID", Path: "PrivateUplinkID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "PrivateUplinkName", Name: "PrivateUplinkName", Path: "PrivateUplinkName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "BandwidthID", Name: "BandwidthID", Path: "BandwidthID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "BandwidthName", Name: "BandwidthName", Path: "BandwidthName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "OperatingSystemID", Name: "OperatingSystemID", Path: "OperatingSystemID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "OperatingSystemFullName", Name: "OperatingSystemFullName", Path: "OperatingSystemFullName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
	}
	f.ChildFields = append(f.ChildFields, childs...)

	return f
}

func RegisterHostDefinition() {
	hostEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: CmdDefaultFields,
		eType:            HostType,
	}
	if err := Registry.Register(hostEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterDedicatedServerDefinition() {
	serverEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", ListHandlerFunc: mapHandler, PageViewHandlerFunc: mapHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			getConfigurationDetailsField(),
		},
		cmdDefaultFields: CmdDefaultFields,
		eType:            DedicatedServerType,
	}
	if err := Registry.Register(serverEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterKubernetesBaremetalNodeDefinition() {
	serverEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", ListHandlerFunc: mapHandler, PageViewHandlerFunc: mapHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			getConfigurationDetailsField(),
		},
		cmdDefaultFields: CmdDefaultFields,
		eType:            KubernetesBaremetalNodeType,
	}
	if err := Registry.Register(serverEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterSBMServerDefinition() {
	serverEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", ListHandlerFunc: mapHandler, PageViewHandlerFunc: mapHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			getConfigurationDetailsField(),
		},
		cmdDefaultFields: CmdDefaultFields,
		eType:            SBMServerType,
	}
	if err := Registry.Register(serverEntity); err != nil {
		log.Fatal(err)
	}
}
