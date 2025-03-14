package entities

import (
	"fmt"
	"io"
	"reflect"
	"slices"
)

// RegistryInterface represents the interface for the EntityRegistry
type RegistryInterface interface {
	Register(entity EntityInterface) error
	GetEntityFromValue(any) EntityInterface
}

// EntityRegistry represents a registry of entities
type EntityRegistry map[reflect.Type]EntityInterface

// EntityInterface represents the interface for an entity
type EntityInterface interface {
	GetType() reflect.Type
	GetDefaultFields() []string
	GetFields() []Field
	Validate(fields []string) error
	SetCmdDefaultFields(cmd string) error
}

// Entity represents the base entity structure
type Entity struct {
	fields           []Field
	cmdDefaultFields map[string][]string
	eType            reflect.Type
}

// HandlerFunc represents a handler function for rendering fields
type HandlerFunc func(w io.Writer, v any, indent string, f *Field) error

// Field represents an entity field
type Field struct {
	ID                  string
	Name                string
	Path                string
	ListHandlerFunc     HandlerFunc
	PageViewHandlerFunc HandlerFunc
	Default             bool
	Parent              *Field
	ChildFields         []Field
}

// PageViewRender renders the field as a page view
func (f *Field) PageViewRender(w io.Writer, v any, indent string) error {
	return f.PageViewHandlerFunc(w, v, indent, f)
}

// ListRender renders the field as a list item
func (f *Field) ListRender(w io.Writer, v any) error {
	return f.ListHandlerFunc(w, v, "", f)
}

// GetName returns the name of the field
func (f *Field) GetName() string {
	return f.Name
}

// GetPath returns the path of the field with its parent path
func (f *Field) GetPath() string {
	if f.Parent != nil {
		return fmt.Sprintf("%s.%s", f.Parent.GetPath(), f.Path)
	}
	return f.Path
}

// Register registers an entity in the registry
func (r *EntityRegistry) Register(entity EntityInterface) error {
	eType := entity.GetType()
	if _, ok := (*r)[eType]; ok {
		return fmt.Errorf("entity already registered: %s", eType)
	}
	(*r)[eType] = entity
	return nil
}

// GetEntityFromValue returns the entity interface from a given value
func (r *EntityRegistry) GetEntityFromValue(v any) (EntityInterface, error) {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if v, ok := (*r)[t]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("entity is not registered: %s", t)
}

// GetType returns the type of the entity
func (e *Entity) GetType() reflect.Type {
	return e.eType
}

// GetDefaultFields returns the default fields for an entity
func (e *Entity) GetDefaultFields() []string {
	result := make([]string, 0)
	for _, field := range e.fields {
		if field.Default {
			result = append(result, field.ID)
		}
	}
	return result
}

// GetAvailableFields returns the available fields for an entity
func (e *Entity) GetFields() []Field {
	return e.fields
}

// Validate checks that all fields match available ones
func (e *Entity) Validate(fields []string) error {
	availableFields := e.GetFields()

	for _, f := range fields {
		found := false
		for _, af := range availableFields {
			if af.ID == f || e.fieldInChildFields(af.ChildFields, f) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("field %s is not found, try --field-list to get available fields", f)
		}
	}

	return nil
}

// SetCmdDefaultFields sets the default fields for an entity based on cmd
func (e *Entity) SetCmdDefaultFields(cmd string) error {
	if defaultFields, ok := e.cmdDefaultFields[cmd]; ok {
		fieldSet := make(map[string]struct{}, len(e.fields))
		for i, f := range e.fields {
			e.fields[i].Default = false
			fieldSet[f.ID] = struct{}{}
		}

		for _, df := range defaultFields {
			if _, exists := fieldSet[df]; !exists {
				return fmt.Errorf("can't find field %s in entity field set", df)
			}
		}

		for i := range e.fields {
			e.fields[i].Default = slices.Contains(defaultFields, e.fields[i].ID)
		}
	}
	return nil
}

// fieldInChildFields checks if a field is in the list of child fields
func (e *Entity) fieldInChildFields(childFields []Field, fieldID string) bool {
	for _, child := range childFields {
		if child.ID == fieldID || e.fieldInChildFields(child.ChildFields, fieldID) {
			return true
		}
	}
	return false
}
