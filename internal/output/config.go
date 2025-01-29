package output

import (
	"fmt"
	"sort"
	"text/tabwriter"
)

type ConfigInfo struct {
	Context  string         `json:"context" yaml:"context"`
	Endpoint string         `json:"endpoint" yaml:"endpoint"`
	Config   map[string]any `json:"config" yaml:"config"`
}

func (f *Formatter) formatConfig(cfgInfo ConfigInfo) error {
	w := tabwriter.NewWriter(f.writer, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "Context:\t%s\n", cfgInfo.Context)
	fmt.Fprintf(w, "Endpoint:\t%s\n", cfgInfo.Endpoint)
	fmt.Fprintln(w, "\nConfiguration:")

	var keys []string
	for k := range cfgInfo.Config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s:\t%v\n", k, cfgInfo.Config[k])
	}

	return w.Flush()
}
