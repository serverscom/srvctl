package entities

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/serverscom/srvctl/internal/output/utils"
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
	GetHeader(Field) string
	GetType() reflect.Type
	GetDefaultFields() []Field
	GetAvailableFields() []Field
	GetFieldsToShow() []Field
	Validate([]string) error
	AddFieldToShow(...Field)
}

// Entity represents the base entity structure
type Entity struct {
	defaultFields []Field
	fieldsToShow  []Field
	eType         reflect.Type
}

// Field represents an entity field
type Field struct {
	Name   string
	Header string
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

// GetHeader returns the header for a given field.
// If Header is not specified, it returns humanized json tag or field name itself.
func (e *Entity) GetHeader(field Field) string {
	if field.Header != "" {
		return field.Header
	}

	var header string

	if f, ok := e.eType.FieldByName(field.Name); ok {
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			header = field.Name
		} else {
			header = strings.Split(jsonTag, ",")[0]
		}
	}
	return utils.Humanize(header)
}

// GetType returns the type of the entity
func (e *Entity) GetType() reflect.Type {
	return e.eType
}

// GetDefaultFields returns the default fields for an entity
func (e *Entity) GetDefaultFields() []Field {
	return e.defaultFields
}

// GetFieldsToShow returns the fields to show for an entity
func (e *Entity) GetFieldsToShow() []Field {
	return e.fieldsToShow
}

// AddFieldToShow adds fields to show for an entity
func (e *Entity) AddFieldToShow(f ...Field) {
	if e.fieldsToShow == nil {
		e.fieldsToShow = make([]Field, 0)
	}
	e.fieldsToShow = append(e.fieldsToShow, f...)
}

// GetAvailableFields returns the available fields for an entity
func (e *Entity) GetAvailableFields() []Field {
	timeType := reflect.TypeOf(time.Time{})
	fields := make([]Field, 0, e.eType.NumField())
	for i := 0; i < e.eType.NumField(); i++ {
		v := e.eType.Field(i)
		t := v.Type

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t.Kind() == reflect.Slice {
			t = t.Elem()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
		}

		// skip nested structs
		if t.Kind() == reflect.Struct && t != timeType {
			continue
		}
		fields = append(fields, Field{Name: v.Name})
	}
	return fields
}

// Validate checks that all fields match available ones
func (e *Entity) Validate(fields []string) error {
	availableFields := e.GetAvailableFields()

	fieldMap := make(map[string]struct{})
	for _, f := range availableFields {
		fieldMap[f.Name] = struct{}{}
	}

	for _, f := range fields {
		if _, exists := fieldMap[f]; !exists {
			return fmt.Errorf("field %s is not found, try --field-list to get available fields", f)
		}
	}

	return nil
}
