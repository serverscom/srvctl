package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	SSHKeyType              = reflect.TypeOf(serverscom.SSHKey{})
	SSHKeyListDefaultFields = []string{"Name", "Fingerprint"}
)

func RegisterSSHKeyDefinition() {
	sshEntity := &Entity{
		fields: []Field{
			{ID: "Name", Name: "Name", Path: "Name", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Fingerprint", Name: "Fingerprint", Path: "Fingerprint", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Labels", Name: "Labels", Path: "Labels", PageViewHandlerFunc: mapPvHandler},
			{ID: "Created", Name: "Created", Path: "Created", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
			{ID: "Updated", Name: "Updated", Path: "Updated", ListHandlerFunc: timeHandler, PageViewHandlerFunc: timeHandler, Default: true},
		},
		cmdDefaultFields: map[string][]string{
			"list": SSHKeyListDefaultFields,
		},
		eType: SSHKeyType,
	}
	if err := Registry.Register(sshEntity); err != nil {
		log.Fatal(err)
	}
}
