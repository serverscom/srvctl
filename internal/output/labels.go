package output

import (
	"fmt"
	"strings"
)

func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "<none>"
	}

	pairs := make([]string, 0, len(labels))
	for k, v := range labels {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, ",")
}
