package output

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/serverscom/srvctl/internal/config"
)

type ContextInfo struct {
	Context  string                 `json:"context" yaml:"context"`
	Endpoint string                 `json:"endpoint" yaml:"endpoint"`
	Config   map[string]interface{} `json:"config" yaml:"config"`
}

func FilterDefaultContexts(contexts []config.Context, defaultContext string, wantDefault bool) []config.Context {
	var filtered []config.Context
	for _, ctx := range contexts {
		isDefault := ctx.Name == defaultContext
		if isDefault == wantDefault {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}

func FormatContexts(contexts []config.Context, defaultContext string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tENDPOINT\tDEFAULT")

	for _, ctx := range contexts {
		isDefault := "*"
		if ctx.Name != defaultContext {
			isDefault = ""
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", ctx.Name, ctx.Endpoint, isDefault)
	}

	return w.Flush()
}

func formatContextInfo(ctx ContextInfo) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "Context:\t%s\n", ctx.Context)
	fmt.Fprintf(w, "Endpoint:\t%s\n", ctx.Endpoint)
	fmt.Fprintln(w, "\nConfiguration:")

	var keys []string
	for k := range ctx.Config {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s:\t%v\n", k, ctx.Config[k])
	}

	return w.Flush()
}
