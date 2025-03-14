package entities

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/serverscom/srvctl/internal/output/utils"
)

const (
	timeFormat = time.RFC3339
)

func stringHandler(w io.Writer, v any, indent string, _ *Field) error {
	// to avoid extra indent between field name and value for nested structs fields
	if indent != "" {
		indent = "\t"
	}
	if v == nil {
		fmt.Fprintf(w, "%s<none>", indent)
		return nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			fmt.Fprintf(w, "%s<none>", indent)
			return nil
		}
		fmt.Fprintf(w, "%s%v", indent, val.Elem().Interface())
		return nil
	}

	fmt.Fprintf(w, "%s%v", indent, v)
	return nil
}

func timeHandler(w io.Writer, v any, indent string, _ *Field) error {
	if indent != "" {
		indent = "\t"
	}
	switch v := v.(type) {
	case time.Time:
		if v.IsZero() {
			_, err := fmt.Fprintf(w, "%s<none>", indent)
			return err
		}
		_, err := fmt.Fprintf(w, "%s%s", indent, v.Format(timeFormat))
		return err
	case *time.Time:
		if v == nil {
			_, err := fmt.Fprintf(w, "%s<none>", indent)
			return err
		}
		_, err := fmt.Fprintf(w, "%s%s", indent, v.Format(timeFormat))
		return err
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func mapHandler(w io.Writer, v interface{}, indent string, _ *Field) error {
	if v == nil {
		_, err := fmt.Fprintf(w, "%s<none>", indent)
		return err
	}

	m, ok := v.(map[string]string)
	if !ok || len(m) == 0 {
		_, err := fmt.Fprintf(w, "%s<empty>", indent)
		return err
	}

	pairs := make([]string, 0, len(m))
	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		pairs = append(pairs, fmt.Sprintf("%s%v=%v", indent, iter.Key(), iter.Value()))
	}

	_, err := fmt.Fprint(w, strings.Join(pairs, "\n"))
	return err
}

func structPVHandler(w io.Writer, v any, indent string, f *Field) error {
	if v == nil {
		fmt.Fprintf(w, "\t<none>")
		return nil
	}

	fmt.Fprintln(w)

	for i, childField := range f.ChildFields {
		fieldValue, err := utils.GetFieldValue(v, childField.Path)
		if err != nil {
			return err
		}

		if childField.PageViewHandlerFunc == nil {
			return fmt.Errorf("no PageViewHandlerFunc defined for field %s", childField.Name)
		}

		fmt.Fprintf(w, "%s%s:", indent, childField.Name)

		if err := childField.PageViewRender(w, fieldValue, indent+"\t"); err != nil {
			return err
		}
		if i < len(f.ChildFields)-1 {
			fmt.Fprintln(w)
		}
	}

	return nil
}
