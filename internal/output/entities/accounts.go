package entities

import (
	"log"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	AccountBalanceType = reflect.TypeOf(serverscom.AccountBalance{})
)

func RegisterAccountDefinition() {
	balanceEntity := &Entity{
		fields: []Field{
			{ID: "CurrentBalance", Name: "CurrentBalance", Path: "CurrentBalance", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "NextInvoiceTotalDue", Name: "NextInvoiceTotalDue", Path: "NextInvoiceTotalDue", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
			{ID: "Currency", Name: "Currency", Path: "Currency", ListHandlerFunc: stringHandler, PageViewHandlerFunc: stringHandler, Default: true},
		},
		eType: AccountBalanceType,
	}
	if err := Registry.Register(balanceEntity); err != nil {
		log.Fatal(err)
	}
}
