package base

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidateFieldPrefixConsistency(t *testing.T) {
	g := NewWithT(t)

	g.Expect(validateFieldPrefixConsistency([]string{"Name", "Fingerprint"})).To(BeNil())
	g.Expect(validateFieldPrefixConsistency([]string{"+Name", "-Fingerprint"})).To(BeNil())
	g.Expect(validateFieldPrefixConsistency([]string{"Name", "+Fingerprint"})).To(HaveOccurred())
	g.Expect(validateFieldPrefixConsistency([]string{"+Name", "Fingerprint"})).To(HaveOccurred())
}
