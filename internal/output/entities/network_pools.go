package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	NetworkPoolType              = reflect.TypeOf(serverscom.NetworkPool{})
	NetworkPoolListDefaultFields = []string{"ID", "Title", "CIDR", "Type"}

	SubnetworkType              = reflect.TypeOf(serverscom.Subnetwork{})
	SubnetworkListDefaultFields = []string{"ID", "Title", "CIDR", "Attached", "InterfaceType"}
)

func RegisterNetworkPoolDefinitions() {
	networkPoolEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CIDR", Name: "CIDR", Path: "CIDR", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationIDs", Name: "LocationIDs", Path: "LocationIDs", PageViewHandlerFunc: slicePvHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},

		cmdDefaultFields: map[string][]string{
			"list": NetworkPoolListDefaultFields,
		},
		eType: NetworkPoolType,
	}
	if err := Registry.Register(networkPoolEntity); err != nil {
		log.Fatal(err)
	}

	subnetworkEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "NetworkPoolID", Name: "NetworkPoolID", Path: "NetworkPoolID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Title", Name: "Title", Path: "Title", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "CIDR", Name: "CIDR", Path: "CIDR", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Attached", Name: "Attached", Path: "Attached", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "InterfaceType", Name: "InterfaceType", Path: "InterfaceType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},

		cmdDefaultFields: map[string][]string{
			"list": SubnetworkListDefaultFields,
		},
		eType: SubnetworkType,
	}
	if err := Registry.Register(subnetworkEntity); err != nil {
		log.Fatal(err)
	}
}
