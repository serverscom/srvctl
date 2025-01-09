package output

import (
	"fmt"
	"text/tabwriter"

	"github.com/serverscom/srvctl/internal/config"
)

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

func (f *Formatter) FormatContexts(contexts []config.Context, defaultContext string) error {
	w := tabwriter.NewWriter(f.writer, 0, 0, 3, ' ', 0)
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
