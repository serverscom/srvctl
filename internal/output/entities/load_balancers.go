package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	LoadBalancerType              = reflect.TypeOf(serverscom.LoadBalancer{})
	L4LoadBalancerType            = reflect.TypeOf(serverscom.L4LoadBalancer{})
	L7LoadBalancerType            = reflect.TypeOf(serverscom.L7LoadBalancer{})
	LoadBalancerListDefaultFields = []string{"ID", "Name", "Type", "Status", "LocationID", "ClusterID"}
)

func RegisterLoadBalancerDefinitions() {
	// Register LoadBalancer
	eLoadBalancer := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ExternalAddresses", Name: "ExternalAddresses", Path: "ExternalAddresses", PageViewHandlerFunc: slicePvHandler},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ClusterID", Name: "ClusterID", Path: "ClusterID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": LoadBalancerListDefaultFields,
		},
		eType: LoadBalancerType,
	}
	if err := Registry.Register(eLoadBalancer); err != nil {
		log.Fatal(err)
	}

	// Register L4LoadBalancer
	eL4LoadBalancer := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ExternalAddresses", Name: "ExternalAddresses", Path: "ExternalAddresses", PageViewHandlerFunc: slicePvHandler},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "StoreLogs", Name: "StoreLogs", Path: "StoreLogs", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "ClusterID", Name: "ClusterID", Path: "ClusterID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": LoadBalancerListDefaultFields,
		},
		eType: L4LoadBalancerType,
	}
	if err := Registry.Register(eL4LoadBalancer); err != nil {
		log.Fatal(err)
	}

	// Register L7LoadBalancer
	eL7LoadBalancer := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "ExternalAddresses", Name: "ExternalAddresses", Path: "ExternalAddresses", PageViewHandlerFunc: slicePvHandler},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Geoip", Name: "Geoip", Path: "Geoip", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "StoreLogs", Name: "StoreLogs", Path: "StoreLogs", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "StoreLogsRegionID", Name: "StoreLogsRegionID", Path: "StoreLogsRegionID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Domains", Name: "Domains", Path: "Domains", PageViewHandlerFunc: slicePvHandler},
			{ID: "ClusterID", Name: "ClusterID", Path: "ClusterID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": LoadBalancerListDefaultFields,
		},
		eType: L7LoadBalancerType,
	}
	if err := Registry.Register(eL7LoadBalancer); err != nil {
		log.Fatal(err)
	}
}
