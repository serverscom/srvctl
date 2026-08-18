package metrics

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []Sample
	}{
		{
			name: "metadata in the format returned by the api",
			input: "# HELP: serverscom_hosts_count Count of the hosts\n" +
				"# TYPE: serverscom_hosts_count gauge\n" +
				`serverscom_hosts_count{location_code="LUX3"} 1` + "\n",
			expected: []Sample{
				{
					Name:   "serverscom_hosts_count",
					Type:   "gauge",
					Help:   "Count of the hosts",
					Labels: map[string]string{"location_code": "LUX3"},
					Value:  1,
				},
			},
		},
		{
			name: "metadata in the standard format",
			input: "# HELP serverscom_hosts_count Count of the hosts\n" +
				"# TYPE serverscom_hosts_count gauge\n" +
				`serverscom_hosts_count{location_code="LUX3"} 1` + "\n",
			expected: []Sample{
				{
					Name:   "serverscom_hosts_count",
					Type:   "gauge",
					Help:   "Count of the hosts",
					Labels: map[string]string{"location_code": "LUX3"},
					Value:  1,
				},
			},
		},
		{
			name:  "escaped label value",
			input: `serverscom_hosts_count{chassis_name="Dell R440 - Silver 4114 - 2.5\"",rack_type="shared"} 2`,
			expected: []Sample{
				{
					Name: "serverscom_hosts_count",
					Labels: map[string]string{
						"chassis_name": `Dell R440 - Silver 4114 - 2.5"`,
						"rack_type":    "shared",
					},
					Value: 2,
				},
			},
		},
		{
			name:  "sample without labels",
			input: "serverscom_hosts_count 42",
			expected: []Sample{
				{Name: "serverscom_hosts_count", Value: 42},
			},
		},
		{
			name:  "sample with empty labels",
			input: "serverscom_hosts_count{} 42",
			expected: []Sample{
				{Name: "serverscom_hosts_count", Labels: map[string]string{}, Value: 42},
			},
		},
		{
			name:  "sample with a timestamp",
			input: `serverscom_hosts_count{location_code="LUX3"} 42 1700000000000`,
			expected: []Sample{
				{
					Name:   "serverscom_hosts_count",
					Labels: map[string]string{"location_code": "LUX3"},
					Value:  42,
				},
			},
		},
		{
			name:  "float and infinite values",
			input: "serverscom_rack_pdu_power_watts 620.5\nserverscom_rack_pdu_current_amperes +Inf",
			expected: []Sample{
				{Name: "serverscom_rack_pdu_power_watts", Value: 620.5},
				{Name: "serverscom_rack_pdu_current_amperes", Value: math.Inf(1)},
			},
		},
		{
			name:     "comments and empty lines only",
			input:    "\n# some comment\n#\n# HELPER not a metadata line\n   \n",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			samples, err := Parse(tc.input)

			g.Expect(err).To(BeNil())
			g.Expect(samples).To(Equal(tc.expected))
		})
	}
}

func TestParseNaN(t *testing.T) {
	g := NewWithT(t)

	samples, err := Parse("serverscom_rack_pdu_power_watts NaN")

	g.Expect(err).To(BeNil())
	g.Expect(samples).To(HaveLen(1))
	g.Expect(math.IsNaN(samples[0].Value)).To(BeTrue())
}

func TestParseErrors(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{name: "no value", input: "serverscom_hosts_count"},
		{name: "not a number", input: "serverscom_hosts_count abc"},
		{name: "extra fields", input: "serverscom_hosts_count 1 2 3"},
		{name: "label without value", input: "serverscom_hosts_count{rack_type} 1"},
		{name: "label value without name", input: `serverscom_hosts_count{="shared"} 1`},
		{name: "unquoted label value", input: "serverscom_hosts_count{rack_type=shared} 1"},
		{name: "unterminated label value", input: `serverscom_hosts_count{rack_type="shared} 1`},
		{name: "unterminated label section", input: `serverscom_hosts_count{rack_type="shared"`},
		{name: "no metric name", input: `{rack_type="shared"} 1`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			samples, err := Parse(tc.input)

			g.Expect(err).To(HaveOccurred())
			g.Expect(samples).To(BeNil())
		})
	}
}
