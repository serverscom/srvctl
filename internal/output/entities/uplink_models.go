package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	UplinkOptionType = reflect.TypeOf(serverscom.UplinkOption{})
)

func RegisterUplinkOptionDefinition() {
	uplinkOptionEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Speed", Name: "Speed", Path: "Speed", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Redundancy", Name: "Redundancy", Path: "Redundancy", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: UplinkOptionType,
	}
	if err := Registry.Register(uplinkOptionEntity); err != nil {
		log.Fatal(err)
	}
}
