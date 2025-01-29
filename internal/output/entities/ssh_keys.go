package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	SSHKeyType = reflect.TypeOf(serverscom.SSHKey{})
)

func RegisterSSHKeyDefinition() {
	sshEntity := &Entity{
		fields: []Field{
			{ID: "Name", Name: "Name", Path: "name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Fingerprint", Name: "Fingerprint", Path: "fingerprint", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "labels", ListHandlerFunc: mapListHandler, PageViewHandlerFunc: mapPageHandler},
			{ID: "Created", Name: "Created", Path: "created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		eType: SSHKeyType,
	}
	if err := Registry.Register(sshEntity); err != nil {
		log.Fatal(err)
	}
}
