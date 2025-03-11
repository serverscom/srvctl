package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/serverscom/srvctl/internal/output/entities"
	"github.com/serverscom/srvctl/internal/output/utils"
)

// ListEntityFields prints available fields
func (f *Formatter) ListEntityFields(fields []entities.Field) {
	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Field ID\tField Name\tSupported Modes")
	for _, field := range fields {
		f.printField(w, field)
	}
	w.Flush()
}

// printField prints a single field to the writer
func (f *Formatter) printField(w *tabwriter.Writer, field entities.Field) {
	modes := []string{}
	if field.ListHandlerFunc != nil {
		modes = append(modes, "list")
	}
	if field.PageViewHandlerFunc != nil {
		modes = append(modes, "page-view")
	}
	if len(modes) > 0 {
		modesStr := strings.Join(modes, ", ")
		fmt.Fprintf(w, "%s\t%s\t%s\n", field.ID, field.Name, modesStr)
	}

	for _, childField := range field.ChildFields {
		f.printField(w, childField)
	}
}

// getOrderedFields returns ordered fields based on the configuration.
// Returns default fields if no fieldsToShow are provided for list mode.
// Returns all available fields if no fieldsToShow are provided for page-view mode.
func (f *Formatter) getOrderedFields(entity entities.EntityInterface) []entities.Field {
	fieldsToShow := f.fieldsToShow

	if len(fieldsToShow) == 0 {
		if f.pageView {
			return entity.GetFields()
		}
		fieldsToShow = entity.GetDefaultFields()
	}

	availableFields := entity.GetFields()
	orderedFields := make([]entities.Field, 0, len(fieldsToShow))

	for _, fieldName := range fieldsToShow {
		if field := f.findFieldByID(availableFields, fieldName); field != nil {
			orderedFields = append(orderedFields, *field)
		}
	}

	return orderedFields
}

// findFieldByID recursively finds a field by its ID in the given fields slice.
func (f *Formatter) findFieldByID(fields []entities.Field, fieldID string) *entities.Field {
	for _, field := range fields {
		if field.ID == fieldID {
			return &field
		}
		if childField := f.findFieldByID(field.ChildFields, fieldID); childField != nil {
			return childField
		}
	}
	return nil
}

// formatText formats the given data as plain text.
func (f *Formatter) formatText(v any) error {
	entity, err := entities.Registry.GetEntityFromValue(v)
	if err != nil {
		return err
	}

	if f.template != nil {
		return f.template.Execute(f.writer, v)
	}

	if f.pageView {
		return f.formatPageView(v, entity)
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 3, ' ', 0)
	defer w.Flush()

	orderedFields := f.getOrderedFields(entity)

	headers := make([]string, 0, len(orderedFields))
	for _, field := range orderedFields {
		headers = append(headers, field.GetName())
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	value := reflect.ValueOf(v)
	return processValue(value, func(item any) error {
		return f.formatRow(w, item, orderedFields)
	})
}

// processValue processes the value of a field and applies any necessary formatting or processing.
func processValue(value reflect.Value, processor func(any) error) error {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Slice:
		for i := range value.Len() {
			if err := processor(value.Index(i).Interface()); err != nil {
				return err
			}
		}
	default:
		return processor(value.Interface())
	}
	return nil
}

// formatRow formats a single row for the given item and fields.
func (f *Formatter) formatRow(w io.Writer, item any, fields []entities.Field) error {
	values := make([]string, 0, len(fields))

	for _, field := range fields {
		fieldValue, err := utils.GetFieldValue(item, field.GetPath())
		if err != nil {
			return err
		}

		var buf strings.Builder

		if field.ListHandlerFunc == nil {
			return fmt.Errorf("no ListHandlerFunc defined for field %s", field.Name)
		}

		if err := field.ListRender(&buf, fieldValue); err != nil {
			return err
		}
		values = append(values, buf.String())
	}

	fmt.Fprintln(w, strings.Join(values, "\t"))
	return nil
}

// formatPageView formats the given data as a page view.
func (f *Formatter) formatPageView(v any, entity entities.EntityInterface) error {
	orderedFields := f.getOrderedFields(entity)
	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	switch value.Kind() {
	case reflect.Slice:
		for i := range value.Len() {
			if err := f.formatPageViewItem(w, value.Index(i).Interface(), orderedFields); err != nil {
				return err
			}
			if i < value.Len()-1 {
				fmt.Fprintln(w, "---")
				w.Flush()
			}
		}
	default:
		return f.formatPageViewItem(w, v, orderedFields)
	}

	return nil
}

// formatPageViewItem formats a single item for the page view mode.
func (f *Formatter) formatPageViewItem(w io.Writer, item any, fields []entities.Field) error {
	for _, field := range fields {
		fieldValue, err := utils.GetFieldValue(item, field.GetPath())
		if err != nil {
			return err
		}

		if field.PageViewHandlerFunc == nil {
			return fmt.Errorf("no PageViewHandlerFunc defined for field %s", field.Name)
		}

		fmt.Fprintf(w, "%s:", field.Name)
		if err := field.PageViewRender(w, fieldValue, "\t"); err != nil {
			return err
		}
		fmt.Fprintln(w)
	}
	return nil
}
