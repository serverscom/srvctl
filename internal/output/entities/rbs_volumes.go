package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	RBSVolumeType = reflect.TypeOf(serverscom.RemoteBlockStorageVolume{})
)

func RegisterRbsVolumeDefinitions() {
	rbsVolumeEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},

		eType: RBSVolumeType,
	}
	if err := Registry.Register(rbsVolumeEntity); err != nil {
		log.Fatal(err)
	}
}
