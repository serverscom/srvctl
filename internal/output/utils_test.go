package output

import (
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/internal/output/entities"
)

func sshKeyEntity(g *WithT) entities.EntityInterface {
	entity, err := entities.Registry.GetEntityFromValue(serverscom.SSHKey{})
	g.Expect(err).To(BeNil())
	return entity
}

func fieldIDs(fields []entities.Field) []string {
	ids := make([]string, 0, len(fields))
	for _, f := range fields {
		ids = append(ids, f.ID)
	}
	return ids
}

func TestGetOrderedFieldsPlain(t *testing.T) {
	g := NewWithT(t)
	entity := sshKeyEntity(g)

	f := &Formatter{cmdName: "list"}
	g.Expect(fieldIDs(f.getOrderedFields(entity))).To(Equal([]string{"Name", "Fingerprint"}))

	f = &Formatter{cmdName: "list", fieldsToShow: []string{"Fingerprint"}}
	g.Expect(fieldIDs(f.getOrderedFields(entity))).To(Equal([]string{"Fingerprint"}))
}

func TestGetOrderedFieldsAdd(t *testing.T) {
	g := NewWithT(t)
	entity := sshKeyEntity(g)

	f := &Formatter{cmdName: "list", fieldsToShow: []string{"+Created"}}
	g.Expect(fieldIDs(f.getOrderedFields(entity))).To(Equal([]string{"Name", "Fingerprint", "Created"}))
}

func TestGetOrderedFieldsRemove(t *testing.T) {
	g := NewWithT(t)
	entity := sshKeyEntity(g)

	f := &Formatter{cmdName: "list", fieldsToShow: []string{"-Fingerprint"}}
	g.Expect(fieldIDs(f.getOrderedFields(entity))).To(Equal([]string{"Name"}))
}

func TestGetOrderedFieldsAddAndRemove(t *testing.T) {
	g := NewWithT(t)
	entity := sshKeyEntity(g)

	f := &Formatter{cmdName: "list", fieldsToShow: []string{"+Created", "-Fingerprint"}}
	g.Expect(fieldIDs(f.getOrderedFields(entity))).To(Equal([]string{"Name", "Created"}))
}

func TestGetOrderedFieldsAddIgnoredInPageView(t *testing.T) {
	g := NewWithT(t)
	entity := sshKeyEntity(g)

	f := &Formatter{cmdName: "list", pageView: true, fieldsToShow: []string{"+Created"}}
	g.Expect(fieldIDs(f.getOrderedFields(entity))).To(Equal(fieldIDs(entity.GetFields())))
}

func TestApplyFieldDeltas(t *testing.T) {
	g := NewWithT(t)

	g.Expect(applyFieldDeltas([]string{"A", "B"}, []string{"+C"})).To(Equal([]string{"A", "B", "C"}))
	g.Expect(applyFieldDeltas([]string{"A", "B"}, []string{"-A"})).To(Equal([]string{"B"}))
	g.Expect(applyFieldDeltas([]string{"A", "B"}, []string{"+B"})).To(Equal([]string{"A", "B"}))
	g.Expect(applyFieldDeltas([]string{"A", "B"}, []string{"-C"})).To(Equal([]string{"A", "B"}))
	g.Expect(applyFieldDeltas([]string{"A", "B"}, []string{"-A", "+A"})).To(Equal([]string{"B", "A"}))
}
