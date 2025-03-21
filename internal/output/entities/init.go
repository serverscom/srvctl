package entities

var (
	Registry = make(EntityRegistry)
)

// init registers all supported entities
func init() {
	RegisterSSHKeyDefinition()
	RegisterHostDefinition()
	RegisterDedicatedServerDefinition()
	RegisterKubernetesBaremetalNodeDefinition()
	RegisterSBMServerDefinition()
	RegisterSSLCertDefinition()
	RegisterSSLCertCustomDefinition()
	RegisterSSLCertLeDefinition()
}
