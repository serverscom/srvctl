package entities

import (
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

func TestHasFieldPrefix(t *testing.T) {
	g := NewWithT(t)

	g.Expect(HasFieldPrefix("+Name")).To(BeTrue())
	g.Expect(HasFieldPrefix("-Name")).To(BeTrue())
	g.Expect(HasFieldPrefix("Name")).To(BeFalse())
	g.Expect(HasFieldPrefix("")).To(BeFalse())
}

func TestStripFieldPrefix(t *testing.T) {
	g := NewWithT(t)

	g.Expect(StripFieldPrefix("+Name")).To(Equal("Name"))
	g.Expect(StripFieldPrefix("-Name")).To(Equal("Name"))
	g.Expect(StripFieldPrefix("Name")).To(Equal("Name"))
}

func TestEntityValidateWithFieldPrefix(t *testing.T) {
	g := NewWithT(t)

	entity, err := Registry.GetEntityFromValue(serverscom.SSHKey{})
	g.Expect(err).To(BeNil())

	g.Expect(entity.Validate([]string{"+Fingerprint"})).To(BeNil())
	g.Expect(entity.Validate([]string{"-Name"})).To(BeNil())
	g.Expect(entity.Validate([]string{"+Unknown"})).To(HaveOccurred())
}
