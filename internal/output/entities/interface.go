package entities

import (
	"fmt"
	"io"
	"reflect"
)

var (
	Registry = make(EntityRegistry)
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
	Validate([]string) error
}

// Entity represents the base entity structure
type Entity struct {
	fields []Field
	eType  reflect.Type
}

type HandlerFunc func(io.Writer, any)

// Field represents an entity field
type Field struct {
	ID                  string
	Name                string
	Path                string
	ListHandlerFunc     HandlerFunc
	PageViewHandlerFunc HandlerFunc
	Default             bool
}

func (f *Field) GetName() string {
	return f.Name
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
			if af.ID == f {
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
