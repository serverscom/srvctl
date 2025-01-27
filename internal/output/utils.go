package output

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/serverscom/srvctl/internal/output/utils"
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

	if f.pageView {
		return f.FormatPageView(v)
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

func (f *Formatter) FormatPageView(data any) error {
	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Slice:
		return f.formatSlice(val, w)
	case reflect.Struct:
		return f.formatStruct(val, w, 0)
	default:
		return fmt.Errorf("unsupported type: %s", val.Kind())
	}
}

func (f *Formatter) formatValue(w io.Writer, name string, fieldValue reflect.Value, field reflect.StructField, indent int) error {
	indentStr := strings.Repeat("\t", indent)

	if !fieldValue.IsValid() || (fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil()) {
		return f.writeField(w, indentStr, name, "<none>")
	}

	switch {
	case fieldValue.Type() == reflect.TypeOf(time.Time{}):
		return f.formatTimeValue(w, indentStr, name, fieldValue)
	case fieldValue.Kind() == reflect.Struct && fieldValue.Type() != reflect.TypeOf(time.Time{}):
		if _, err := fmt.Fprintf(w, "%s%s:\n", indentStr, name); err != nil {
			return err
		}
		return f.formatStruct(fieldValue, w, indent+1)
	case fieldValue.Kind() == reflect.Slice:
		return f.formatSliceValue(w, indentStr, name, fieldValue, field, indent)
	case fieldValue.Kind() == reflect.Map:
		return f.formatMapValue(w, indentStr, name, fieldValue)
	default:
		return f.writeField(w, indentStr, name, fieldValue.Interface())
	}
}

func (f *Formatter) formatTimeValue(w io.Writer, indentStr, name string, fieldValue reflect.Value) error {
	t := fieldValue.Interface().(time.Time)
	if t.IsZero() {
		return f.writeField(w, indentStr, name, "<none>")
	}
	return f.writeField(w, indentStr, name, t.Format("2006-01-02 15:04:05"))
}

func (f *Formatter) formatSliceValue(w io.Writer, indentStr, name string, fieldValue reflect.Value, field reflect.StructField, indent int) error {
	if fieldValue.IsNil() {
		return f.writeField(w, indentStr, name, "<none>")
	}
	if fieldValue.Len() == 0 {
		return f.writeField(w, indentStr, name, "<empty>")
	}

	if fieldValue.Type().Elem().Kind() == reflect.Struct {
		if _, err := fmt.Fprintf(w, "%s%s:\n", indentStr, name); err != nil {
			return err
		}
		return f.formatStructSlice(w, fieldValue, field, indent)
	}

	elements := make([]string, fieldValue.Len())
	for i := 0; i < fieldValue.Len(); i++ {
		elements[i] = fmt.Sprint(fieldValue.Index(i).Interface())
	}
	return f.writeField(w, indentStr, name, strings.Join(elements, ", "))
}

func (f *Formatter) formatStructSlice(w io.Writer, fieldValue reflect.Value, field reflect.StructField, indent int) error {
	if h, ok := f.rendererHandlers[field.Type]; ok {
		return h(fieldValue.Interface(), f, indent+1)
	}

	indentStr := strings.Repeat("\t", indent+1)
	for i := 0; i < fieldValue.Len(); i++ {
		elem := fieldValue.Index(i)
		key := getNameOrID(elem)
		if key == "" {
			return fmt.Errorf("can't find custom renderer or name/id field for %s", field.Type)
		}
		if _, err := fmt.Fprintf(w, "%s%s:\n", indentStr, key); err != nil {
			return err
		}
		if err := f.formatStruct(elem, w, indent+2); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) formatMapValue(w io.Writer, indentStr, name string, fieldValue reflect.Value) error {
	if fieldValue.IsNil() {
		return f.writeField(w, indentStr, name, "<none>")
	}
	if fieldValue.Len() == 0 {
		return f.writeField(w, indentStr, name, "<empty>")
	}

	pairs := make([]string, 0, fieldValue.Len())
	iter := fieldValue.MapRange()
	for iter.Next() {
		pairs = append(pairs, fmt.Sprintf("%v=%v", iter.Key(), iter.Value()))
	}
	return f.writeField(w, indentStr, name, strings.Join(pairs, "\n"+indentStr+"\t"))
}

func (f *Formatter) writeField(w io.Writer, indent, name string, value interface{}) error {
	_, err := fmt.Fprintf(w, "%s%s:\t%v\n", indent, name, value)
	return err
}

func (f *Formatter) formatSlice(val reflect.Value, w io.Writer) error {
	for i := 0; i < val.Len(); i++ {
		if err := f.formatStruct(val.Index(i), w, 0); err != nil {
			return err
		}
		if i < val.Len()-1 {
			if _, err := fmt.Fprintf(w, "---\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *Formatter) formatStruct(val reflect.Value, w io.Writer, indent int) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got: %s", val.Kind())
	}

	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		if !f.shouldShowField(field.Name) {
			continue
		}

		fieldValue := val.Field(i)
		if fieldValue.Kind() == reflect.Pointer {
			fieldValue = fieldValue.Elem()
		}

		name := f.getFieldName(field)
		if err := f.formatValue(w, name, fieldValue, field, indent); err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) shouldShowField(fieldName string) bool {
	return len(f.fieldsToShow) == 0 || slices.Contains(f.fieldsToShow, fieldName)
}

func (f *Formatter) getFieldName(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		return utils.Humanize(jsonTag)
	}
	return field.Name
}

func getNameOrID(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}

	var priorityFields = []string{"name", "id"}

	for _, priorityField := range priorityFields {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if strings.EqualFold(field.Name, priorityField) {
				return fmt.Sprint(v.Field(i).Interface())
			}
		}
	}

	return ""
}
