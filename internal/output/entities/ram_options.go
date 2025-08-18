package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	RAMOptionType = reflect.TypeOf(serverscom.RAMOption{})
)

func RegisterRAMOptionDefinition() {
	ramOptionEntity := &Entity{
		fields: []Field{
			{ID: "RAM", Name: "RAM", Path: "RAM", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Type", Name: "Type", Path: "Type", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},

		eType: RAMOptionType,
	}
	if err := Registry.Register(ramOptionEntity); err != nil {
		log.Fatal(err)
	}
}
