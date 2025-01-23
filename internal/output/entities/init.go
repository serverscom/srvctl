package entities

// init registers all supported entities
// to use custom header in table output:
//
//	Entity{
//		Fields: []Field{
//			{Name: "Name"},
//			{Name: "Fingerprint", Header: "CUSTOM FINGERPRINT"},
//			{Name: "Created"},
//			{Name: "Updated", Header: "UPDATED_AT"},
//		},
//	}
func init() {
	RegisterSSHKeyDefinition()
}
