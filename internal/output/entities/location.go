package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	LocationType = reflect.TypeOf(serverscom.Location{})
)

func RegisterLocationDefinition() {
	locationEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Code", Name: "Code", Path: "Code", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "SupportedFeatures", Name: "SupportedFeatures", Path: "SupportedFeatures", PageViewHandlerFunc: slicePvHandler},
			{ID: "L2SegmentsEnabled", Name: "L2SegmentsEnabled", Path: "L2SegmentsEnabled", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "PrivateRacksEnabled", Name: "PrivateRacksEnabled", Path: "PrivateRacksEnabled", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "LoadBalancersEnabled", Name: "LoadBalancersEnabled", Path: "LoadBalancersEnabled", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
		},
		eType: LocationType,
	}
	if err := Registry.Register(locationEntity); err != nil {
		log.Fatal(err)
	}
}
