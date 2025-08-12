package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	BandwidthOptionType = reflect.TypeOf(serverscom.BandwidthOption{})
)

func RegisterBandwidthOptionDefinition() {
	bandwidthOptionEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Commit", Name: "Commit", Path: "Commit", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: BandwidthOptionType,
	}

	if err := Registry.Register(bandwidthOptionEntity); err != nil {
		log.Fatal(err)
	}
}
