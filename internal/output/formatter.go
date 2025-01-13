package output

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"gopkg.in/yaml.v3"
)

// Formatter represents formatter struct with custom io.Writer
type Formatter struct {
	writer io.Writer
}

// NewFormatter creates new formatter with specified io.Writer
func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{writer: w}
}

// Format formats data according to format
func (f *Formatter) Format(v interface{}, format string) error {
	// handle single element slices in JSON and YAML formats
	if item, ok := isOneElementSlice(v); ok && (format == "json" || format == "yaml") {
		v = item
	}
	return f.format(v, format)
}

// FormatList formats list of elements according to format
func (f *Formatter) FormatList(v interface{}, format string) error {
	return f.format(v, format)
}

// format formats data according to format
func (f *Formatter) format(v interface{}, format string) error {
	switch format {
	case "json":
		data, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			return err
		}
		_, err = f.writer.Write(data)
		return err
	case "yaml":
		return yaml.NewEncoder(f.writer).Encode(v)
	default:
		return f.formatText(v)
	}
}

// formatText formats data in text format
func (f *Formatter) formatText(v interface{}) error {
	switch data := v.(type) {
	case ConfigInfo:
		return f.formatConfig(data)
	case []serverscom.SSHKey:
		return f.formatSSHKeys(data)
	default:
		return fmt.Errorf("unsupported type for text output: %T", v)
	}
}

func isOneElementSlice(v interface{}) (interface{}, bool) {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Slice {
		return nil, false
	}

	if val.Len() != 1 {
		return nil, false
	}

	return val.Index(0).Interface(), true
}
