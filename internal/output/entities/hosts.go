package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	HostType                                 = reflect.TypeOf(serverscom.Host{})
	DedicatedServerType                      = reflect.TypeOf(serverscom.DedicatedServer{})
	KubernetesBaremetalNodeType              = reflect.TypeOf(serverscom.KubernetesBaremetalNode{})
	SBMServerType                            = reflect.TypeOf(serverscom.SBMServer{})
	HostConnectionType                       = reflect.TypeOf(serverscom.HostConnection{})
	HostPowerFeedType                        = reflect.TypeOf(serverscom.HostPowerFeed{})
	HostDriveSlotType                        = reflect.TypeOf(serverscom.HostDriveSlot{})
	HostPTRRecordType                        = reflect.TypeOf(serverscom.PTRRecord{})
	HostNetworkType                          = reflect.TypeOf(serverscom.Network{})
	HostFeatureType                          = reflect.TypeOf(serverscom.DedicatedServerFeature{})
	HostServiceType                          = reflect.TypeOf(serverscom.DedicatedServerService{})
	HostOOBCredsType                         = reflect.TypeOf(serverscom.DedicatedServerOOBCredentials{})
	HostListDefaultFields                    = []string{"ID", "Type", "Title", "LocationCode", "Status", "PublicIPv4Address"}
	DedicatedServerListDefaultFields         = []string{"ID", "Title", "RackID", "LocationCode", "Status", "PublicIPv4Address"}
	KubernetesBaremetalNodeListDefaultFields = []string{"ID", "KubernetesClusterNodeNumber", "Title", "LocationCode", "Status", "PublicIPv4Address"}
	SBMServerListDefaultFields               = []string{"ID", "Title", "RackID", "LocationCode", "Status", "PublicIPv4Address"}

	NetworkListDefaultFields = []string{"ID", "Title", "Status", "CIDR", "Family", "InterfaceType", "DistributionMethod", "Additional"}
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

func getDriveModel() Field {
	f := Field{
		ID:                  "DriveModel",
		Name:                "DriveModel",
		Path:                "DriveModel",
		PageViewHandlerFunc: structPVHandler,
	}
	childs := []Field{
		{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "Capacity", Name: "Capacity", Path: "Capacity", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "Interface", Name: "Interface", Path: "Interface", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "FormFactor", Name: "FormFactor", Path: "FormFactor", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
		{ID: "MediaType", Name: "MediaType", Path: "MediaType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Parent: &f},
	}
	f.ChildFields = append(f.ChildFields, childs...)

	return f
}

func RegisterHostDefinition() {
	hostEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			// {ID: "RackID", Name: "RackID", Path: "RackID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			// {ID: "LeaseStart", Name: "LeaseStart", Path: "LeaseStart", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			// {ID: "OobIPv4Address", Name: "OobIPv4Address", Path: "OobIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": HostListDefaultFields,
		},
		eType: HostType,
	}
	if err := Registry.Register(hostEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterDedicatedServerDefinition() {
	serverEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RackID", Name: "RackID", Path: "RackID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LeaseStart", Name: "LeaseStart", Path: "LeaseStart", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "OobIPv4Address", Name: "OobIPv4Address", Path: "OobIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			getConfigurationDetailsField(),
		},
		cmdDefaultFields: map[string][]string{
			"list": DedicatedServerListDefaultFields,
		},
		eType: DedicatedServerType,
	}
	if err := Registry.Register(serverEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterKubernetesBaremetalNodeDefinition() {
	serverEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "KubernetesClusterId", Name: "KubernetesClusterId", Path: "KubernetesClusterId", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "KubernetesClusterNodeId", Name: "KubernetesClusterNodeId", Path: "KubernetesClusterNodeId", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "KubernetesClusterNodeNumber", Name: "KubernetesClusterNodeNumber", Path: "KubernetesClusterNodeNumber", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "RackID", Name: "RackID", Path: "RackID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LeaseStart", Name: "LeaseStart", Path: "LeaseStart", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "OobIPv4Address", Name: "OobIPv4Address", Path: "OobIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			getConfigurationDetailsField(),
		},
		cmdDefaultFields: map[string][]string{
			"list": KubernetesBaremetalNodeListDefaultFields,
		},
		eType: KubernetesBaremetalNodeType,
	}
	if err := Registry.Register(serverEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterSBMServerDefinition() {
	serverEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RackID", Name: "RackID", Path: "RackID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "OperationalStatus", Name: "OperationalStatus", Path: "OperationalStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PowerStatus", Name: "PowerStatus", Path: "PowerStatus", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Configuration", Name: "Configuration", Path: "Configuration", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateIPv4Address", Name: "PrivateIPv4Address", Path: "PrivateIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PublicIPv4Address", Name: "PublicIPv4Address", Path: "PublicIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LeaseStart", Name: "LeaseStart", Path: "LeaseStart", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ScheduledRelease", Name: "ScheduledRelease", Path: "ScheduledRelease", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler},
			{ID: "OobIPv4Address", Name: "OobIPv4Address", Path: "OobIPv4Address", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			getConfigurationDetailsField(),
		},
		cmdDefaultFields: map[string][]string{
			"list": SBMServerListDefaultFields,
		},
		eType: SBMServerType,
	}
	if err := Registry.Register(serverEntity); err != nil {
		log.Fatal(err)
	}
}

func RegisterHostsSubDefinitions() {
	hostConnectionEntity := &Entity{
		fields: []Field{
			{ID: "Port", Name: "Port", Path: "Port", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "MACAddress", Name: "MACAddress", Path: "MACAddress", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: HostConnectionType,
	}
	if err := Registry.Register(hostConnectionEntity); err != nil {
		log.Fatal(err)
	}

	hostPowerFeedEntity := &Entity{
		fields: []Field{
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: HostPowerFeedType,
	}
	if err := Registry.Register(hostPowerFeedEntity); err != nil {
		log.Fatal(err)
	}

	hostDriveSlotEntity := &Entity{
		fields: []Field{
			{ID: "Position", Name: "Position", Path: "Position", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Interface", Name: "Interface", Path: "Interface", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "FormFactor", Name: "FormFactor", Path: "FormFactor", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			getDriveModel(),
		},
		eType: HostDriveSlotType,
	}
	if err := Registry.Register(hostDriveSlotEntity); err != nil {
		log.Fatal(err)
	}

	hostPTRRecordEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "IP", Name: "IP", Path: "IP", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Domain", Name: "Domain", Path: "Domain", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Priority", Name: "Priority", Path: "Priority", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "TTL", Name: "TTL", Path: "TTL", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: HostPTRRecordType,
	}
	if err := Registry.Register(hostPTRRecordEntity); err != nil {
		log.Fatal(err)
	}

	hostNetworkEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CIDR", Name: "CIDR", Path: "Cidr", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Family", Name: "Family", Path: "Family", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "InterfaceType", Name: "InterfaceType", Path: "InterfaceType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "DistributionMethod", Name: "DistributionMethod", Path: "DistributionMethod", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Additional", Name: "Additional", Path: "Additional", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		eType: HostNetworkType,
		cmdDefaultFields: map[string][]string{
			"list-networks": NetworkListDefaultFields,
		},
	}
	if err := Registry.Register(hostNetworkEntity); err != nil {
		log.Fatal(err)
	}

	hostFeatureEntity := &Entity{
		fields: []Field{
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: HostFeatureType,
	}
	if err := Registry.Register(hostFeatureEntity); err != nil {
		log.Fatal(err)
	}

	hostServiceEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Currency", Name: "Currency", Path: "Currency", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Label", Name: "Label", Path: "Label", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "StartedAt", Name: "StartedAt", Path: "StartedAt", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "FinishedAt", Name: "FinishedAt", Path: "FinishedAt", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "DateFrom", Name: "DateFrom", Path: "DateFrom", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "DateTo", Name: "DateTo", Path: "DateTo", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "UsageQuantity", Name: "UsageQuantity", Path: "UsageQuantity", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Tax", Name: "Tax", Path: "Tax", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Total", Name: "Total", Path: "Total", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Subtotal", Name: "Subtotal", Path: "Subtotal", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "DiscountRate", Name: "DiscountRate", Path: "DiscountRate", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
		},
		eType: HostServiceType,
	}
	if err := Registry.Register(hostServiceEntity); err != nil {
		log.Fatal(err)
	}

	hostOOBCredsEntity := &Entity{
		fields: []Field{
			{ID: "Login", Name: "Login", Path: "Login", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Secret", Name: "Secret", Path: "Secret", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: HostOOBCredsType,
	}
	if err := Registry.Register(hostOOBCredsEntity); err != nil {
		log.Fatal(err)
	}
}
