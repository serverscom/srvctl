package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	CloudVolumeType              = reflect.TypeOf(serverscom.CloudBlockStorageVolume{})
	CloudVolumeListDefaultFields = []string{"ID", "Name", "RegionID", "Size", "Description"}
)

func RegisterCloudVolumeDefinition() {
	cloudVolumeEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RegionID", Name: "Region ID", Path: "RegionID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Size", Name: "Size", Path: "Size", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Description", Name: "Description", Path: "Description", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": CloudVolumeListDefaultFields,
		},
		eType: CloudVolumeType,
	}
	if err := Registry.Register(cloudVolumeEntity); err != nil {
		log.Fatal(err)
	}
}
