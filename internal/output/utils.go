package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/serverscom/srvctl/internal/output/entities"
)

// ListEntityFields prints fields
func (f *Formatter) ListEntityFields(fields []entities.Field) {
	for _, field := range fields {
		fmt.Fprintln(f.writer, field.Name)
	}
}

// getFieldValue returns the value of a struct field
func getFieldValue(v any, fieldName string) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return "", fmt.Errorf("field %s not found", fieldName)
	}

	if !field.IsValid() || (field.Kind() == reflect.Ptr && field.IsNil()) {
		return "<none>", nil
	}

	fieldType := field.Type()

	switch {
	case fieldType == reflect.TypeOf(time.Time{}):
		t := field.Interface().(time.Time)
		if t.IsZero() {
			return "<none>", nil
		}
		return t.Format("2006-01-02 15:04:05"), nil

	case fieldType.Kind() == reflect.Map:
		if field.IsNil() {
			return "<none>", nil
		}
		if field.Len() == 0 {
			return "<empty>", nil
		}
		pairs := make([]string, 0, field.Len())
		iter := field.MapRange()
		for iter.Next() {
			pairs = append(pairs, fmt.Sprintf("%v=%v", iter.Key(), iter.Value()))
		}
		return strings.Join(pairs, ","), nil

	case fieldType.Kind() == reflect.Slice:
		if field.IsNil() {
			return "<none>", nil
		}
		if field.Len() == 0 {
			return "<empty>", nil
		}
		elements := make([]string, field.Len())
		for i := 0; i < field.Len(); i++ {
			elements[i] = fmt.Sprint(field.Index(i).Interface())
		}
		return strings.Join(elements, ","), nil
	}

	return fmt.Sprintf("%v", field.Interface()), nil
}

// formatText formats the given entity as text.
// Supports only registered entities.
// If no fields are specified, it shows entity default fields.
func (f *Formatter) formatText(v any) error {
	entity, err := entities.Registry.GetEntityFromValue(v)
	if err != nil {
		return err
	}

	if f.template != nil {
		return f.template.Execute(f.writer, v)
	}

	defaultFields := entity.GetDefaultFields()

	if len(f.fieldsToShow) > 0 {
		for _, fieldName := range f.fieldsToShow {
			found := false
			for _, field := range defaultFields {
				if field.Name == fieldName {
					entity.AddFieldToShow(field)
					found = true
					break
				}
			}
			if !found {
				entity.AddFieldToShow(entities.Field{Name: fieldName})
			}
		}
	} else {
		entity.AddFieldToShow(defaultFields...)
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 3, ' ', 0)
	return printTable(w, v, entity)
}

// printTable prints the table with headers and values.
func printTable(w *tabwriter.Writer, v any, entity entities.EntityInterface) error {
	fieldsToShow := entity.GetFieldsToShow()
	headers := make([]string, 0, len(fieldsToShow))
	for _, field := range fieldsToShow {
		h := entity.GetHeader(field)
		headers = append(headers, strings.ToUpper(h))
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	if err := printFields(w, v, fieldsToShow); err != nil {
		return err
	}

	return w.Flush()
}

// printFields prints the values for each field in the slice or single object
func printFields(w io.Writer, v any, fieldsToShow []entities.Field) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			if err := printRow(w, val.Index(i).Interface(), fieldsToShow); err != nil {
				return err
			}
		}
		return nil
	}

	return printRow(w, v, fieldsToShow)
}

// printRow prints a single row of the table with values for each field
func printRow(w io.Writer, v any, fieldsToShow []entities.Field) error {
	values := make([]string, 0, len(fieldsToShow))
	for _, field := range fieldsToShow {
		value, err := getFieldValue(v, field.Name)
		if err != nil {
			return err
		}
		values = append(values, value)
	}
	fmt.Fprintln(w, strings.Join(values, "\t"))
	return nil
}
