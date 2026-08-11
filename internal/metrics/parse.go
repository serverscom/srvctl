// Package metrics parses the Prometheus text exposition format returned by the
// metrics endpoints and folds it into rows suitable for tabular output.
package metrics

import (
	"fmt"
	"strconv"
	"strings"
)

// Sample represents a single metric sample.
type Sample struct {
	Name   string
	Type   string
	Help   string
	Labels map[string]string
	Value  float64
}

// Parse parses metrics in the Prometheus text exposition format.
func Parse(s string) ([]Sample, error) {
	var samples []Sample
	help := make(map[string]string)
	types := make(map[string]string)

	for i, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			parseComment(line, help, types)
			continue
		}

		sample, err := parseSample(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		samples = append(samples, sample)
	}

	for i := range samples {
		samples[i].Help = help[samples[i].Name]
		samples[i].Type = types[samples[i].Name]
	}

	return samples, nil
}

// parseComment fills help and type metadata from a comment line.
// The API emits "# HELP: <name> <help>", the standard format is "# HELP <name> <help>",
// both are accepted. Any other comment is ignored.
func parseComment(line string, help, types map[string]string) {
	rest := strings.TrimSpace(strings.TrimPrefix(line, "#"))

	var target map[string]string
	switch {
	case strings.HasPrefix(rest, "HELP"):
		rest, target = rest[len("HELP"):], help
	case strings.HasPrefix(rest, "TYPE"):
		rest, target = rest[len("TYPE"):], types
	default:
		return
	}
	if rest == "" || (rest[0] != ':' && rest[0] != ' ' && rest[0] != '\t') {
		return
	}

	rest = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(rest), ":"))
	name, value, _ := strings.Cut(rest, " ")
	if name == "" {
		return
	}
	target[name] = strings.TrimSpace(value)
}

// parseSample parses a single sample line, e.g. name{label="value"} 42.
func parseSample(line string) (Sample, error) {
	i := strings.IndexAny(line, "{ \t")
	if i <= 0 {
		return Sample{}, fmt.Errorf("malformed sample %q", line)
	}

	sample := Sample{Name: line[:i]}
	rest := line[i:]

	if strings.HasPrefix(rest, "{") {
		labels, remainder, err := parseLabels(rest)
		if err != nil {
			return Sample{}, err
		}
		sample.Labels = labels
		rest = remainder
	}

	// a sample value can be followed by an optional timestamp
	fields := strings.Fields(rest)
	if len(fields) == 0 || len(fields) > 2 {
		return Sample{}, fmt.Errorf("malformed sample %q: expected a value, got %d fields", line, len(fields))
	}

	value, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return Sample{}, fmt.Errorf("malformed sample %q: %w", line, err)
	}
	sample.Value = value

	return sample, nil
}

// parseLabels parses a label section and returns the labels along with
// everything that follows the closing brace.
func parseLabels(s string) (map[string]string, string, error) {
	labels := make(map[string]string)

	i := 1 // skip the opening brace
	for {
		for i < len(s) && (s[i] == ',' || s[i] == ' ' || s[i] == '\t') {
			i++
		}
		if i >= len(s) {
			return nil, "", fmt.Errorf("unterminated label section in %q", s)
		}
		if s[i] == '}' {
			return labels, s[i+1:], nil
		}

		start := i
		for i < len(s) && s[i] != '=' && s[i] != '}' {
			i++
		}
		name := strings.TrimSpace(s[start:i])
		if i >= len(s) || s[i] != '=' || name == "" {
			return nil, "", fmt.Errorf("malformed label name in %q", s)
		}
		i++

		if i >= len(s) || s[i] != '"' {
			return nil, "", fmt.Errorf("expected a quoted value for label %q in %q", name, s)
		}
		i++

		value, next, err := parseLabelValue(s, i)
		if err != nil {
			return nil, "", err
		}
		labels[name] = value
		i = next
	}
}

// parseLabelValue reads a quoted label value starting at i and returns it
// unescaped along with the position right after the closing quote.
func parseLabelValue(s string, i int) (string, int, error) {
	var value strings.Builder

	for i < len(s) {
		switch s[i] {
		case '"':
			return value.String(), i + 1, nil
		case '\\':
			i++
			if i >= len(s) {
				break
			}
			switch s[i] {
			case 'n':
				value.WriteByte('\n')
			case 't':
				value.WriteByte('\t')
			default:
				value.WriteByte(s[i])
			}
			i++
		default:
			value.WriteByte(s[i])
			i++
		}
	}

	return "", 0, fmt.Errorf("unterminated label value in %q", s)
}
