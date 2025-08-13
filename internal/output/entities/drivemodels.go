package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	DriveModelOptionType = reflect.TypeOf(serverscom.DriveModel{})
)

func RegisterDriveModelOptionDefinition() {
	driveModelOptionEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Capacity", Name: "Capacity", Path: "Capacity", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Interface", Name: "Interface", Path: "Interface", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "FormFactor", Name: "FormFactor", Path: "FormFactor", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "MediaType", Name: "MediaType", Path: "MediaType", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: DriveModelOptionType,
	}
	if err := Registry.Register(driveModelOptionEntity); err != nil {
		log.Fatal(err)
	}
}
