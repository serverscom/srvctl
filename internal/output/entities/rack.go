package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	RackType = reflect.TypeOf(serverscom.Rack{})
)

func RegisterRackDefinition() {
	rackEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationID", Name: "LocationID", Path: "LocationID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "LocationCode", Name: "LocationCode", Path: "LocationCode", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
		},
		eType: RackType,
	}
	if err := Registry.Register(rackEntity); err != nil {
		log.Fatal(err)
	}
}
