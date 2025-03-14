package output

import (
	"encoding/json"
	"html/template"
	"io"

	"github.com/serverscom/srvctl/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Formatter represents formatter struct with custom io.Writer
type Formatter struct {
	writer       io.Writer
	output       string
	templateStr  string
	template     *template.Template
	pageView     bool
	fieldsToShow []string
	fieldList    bool
	header       bool
}

// NewFormatter creates new formatter with specified io.Writer
func NewFormatter(cmd *cobra.Command, manager *config.Manager) *Formatter {
	output, _ := manager.GetResolvedStringValue(cmd, "output")
	pageView, _ := manager.GetResolvedBoolValue(cmd, "page-view")
	template, _ := manager.GetResolvedStringValue(cmd, "template")
	fields, _ := manager.GetResolvedStringSliceValue(cmd, "field")
	fieldList, _ := manager.GetResolvedBoolValue(cmd, "field-list")
	noHeader, _ := manager.GetResolvedBoolValue(cmd, "no-header")

	return &Formatter{
		writer:       cmd.OutOrStdout(),
		output:       output,
		templateStr:  template,
		pageView:     pageView,
		fieldsToShow: fields,
		fieldList:    fieldList,
		header:       !noHeader,
	}
}

// GetOutput returns output type
func (f *Formatter) GetOutput() string {
	return f.output
}

// GetTemplateStr returns template string
func (f *Formatter) GetTemplateStr() string {
	return f.templateStr
}

// SetTemplate sets template
func (f *Formatter) SetTemplate(t *template.Template) {
	f.template = t
}

// Format formats data according to format
func (f *Formatter) Format(v any) error {
	switch f.output {
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
		return f.FormatText(v)
	}
}

// formatText formats data in text format
func (f *Formatter) FormatText(v any) error {
	switch data := v.(type) {
	case ConfigInfo:
		return f.formatConfig(data)
	default:
		return f.formatText(data)
	}
}
