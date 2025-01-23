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
		defaultFields: []Field{
			{Name: "Name"},
			{Name: "Fingerprint"},
			{Name: "Created"},
			{Name: "Updated"},
		},
		eType: SSHKeyType,
	}
	if err := Registry.Register(sshEntity); err != nil {
		log.Fatal(err)
	}
}
