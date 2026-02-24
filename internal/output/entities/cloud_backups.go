package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	CloudBlockStorageBackupType              = reflect.TypeOf(serverscom.CloudBlockStorageBackup{})
	CloudBlockStorageBackupListDefaultFields = []string{"ID", "Name", "Status", "Size"}
)

func RegisterCloudBackupDefinition() {
	backupEntity := &Entity{
		fields: []Field{
			{ID: "ID", Name: "ID", Path: "ID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Status", Name: "Status", Path: "Status", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Size", Name: "Size", Path: "Size", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "RegionID", Name: "Region ID", Path: "RegionID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "OpenstackUUID", Name: "Openstack UUID", Path: "OpenstackUUID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "OpenstackVolumeUUID", Name: "Openstack Volume UUID", Path: "OpenstackVolumeUUID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": CloudBlockStorageBackupListDefaultFields,
		},
		eType: CloudBlockStorageBackupType,
	}
	if err := Registry.Register(backupEntity); err != nil {
		log.Fatal(err)
	}
}
