package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	OperatingSystemOptionType = reflect.TypeOf(serverscom.OperatingSystemOption{})
)

func RegisterOperatingSystemOptionDefinition() {
	osOptionEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "FullName", Name: "Full Name", Path: "FullName", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Version", Name: "Version", Path: "Version", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Arch", Name: "Arch", Path: "Arch", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Filesystems", Name: "Filesystems", Path: "Filesystems", PageViewHandlerFunc: slicePvHandler},
		},

		eType: OperatingSystemOptionType,
	}
	if err := Registry.Register(osOptionEntity); err != nil {
		log.Fatal(err)
	}
}
