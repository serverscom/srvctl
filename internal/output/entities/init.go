package entities

var (
	Registry = make(EntityRegistry)
)

// init registers all supported entities
func init() {
	RegisterSSHKeyDefinition()
	RegisterHostDefinition()
	RegisterHostsSubDefinitions()
	RegisterDedicatedServerDefinition()
	RegisterKubernetesBaremetalNodeDefinition()
	RegisterSBMServerDefinition()
	RegisterSSLCertDefinition()
	RegisterSSLCertCustomDefinition()
	RegisterSSLCertLeDefinition()
	RegisterLoadBalancerDefinitions()
	RegisterRackDefinition()
	RegisterInvoiceDefinition()
	RegisterAccountDefinition()
	RegisterLocationDefinition()
	RegisterKubernetesClusterDefinition()
	RegisterKubernetesClusterNodeDefinition()
}
