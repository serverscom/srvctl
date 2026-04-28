package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	RBSVolumeCredentialsType = reflect.TypeOf(serverscom.RemoteBlockStorageVolumeCredentials{})
)

func RegisterRbsVolumeCredentialsDefinition() {
	rbsCredentialsEntity := &Entity{
		fields: []Field{
			{ID: "VolumeID", Name: "Volume ID", Path: "VolumeID", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Username", Name: "Username", Path: "Username", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Password", Name: "Password", Path: "Password", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "TargetIQN", Name: "Target IQN", Path: "TargetIQN", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "IPAddress", Name: "IP Address", Path: "IPAddress", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: RBSVolumeCredentialsType,
	}
	if err := Registry.Register(rbsCredentialsEntity); err != nil {
		log.Fatal(err)
	}
}
