package metrics

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "metrics")

func readFixture(name string) string {
	return string(testutils.ReadFixture(filepath.Join(fixtureBasePath, name)))
}

type metricsTestCase struct {
	name           string
	args           []string
	metrics        string
	expectedOutput string
	expectedErrOut string
	// noAPICall is set for cases failing before the API is called
	noAPICall bool
	// apiError makes the API call fail
	apiError    bool
	expectError bool
}

func TestHostsCmd(t *testing.T) {
	hostsMetrics := readFixture("hosts_input.txt")

	testCases := []metricsTestCase{
		{
			name:           "get hosts metrics in default format",
			metrics:        hostsMetrics,
			expectedOutput: readFixture("hosts.txt"),
		},
		{
			name:           "get hosts metrics in page view",
			args:           []string{"--page-view"},
			metrics:        hostsMetrics,
			expectedOutput: readFixture("hosts_page_view.txt"),
		},
		{
			name:           "get hosts metrics without header",
			args:           []string{"--no-header"},
			metrics:        hostsMetrics,
			expectedOutput: readFixture("hosts_no_header.txt"),
		},
		{
			name:           "get hosts metrics with fields",
			args:           []string{"-f", "HostID", "-f", "ChassisName", "-f", "TotalSent"},
			metrics:        hostsMetrics,
			expectedOutput: readFixture("hosts_field.txt"),
		},
		{
			name:           "get hosts metrics with template",
			args:           []string{"-t", `{{range .}}{{.HostID}} {{.TotalSent}}\n{{end}}`},
			metrics:        hostsMetrics,
			expectedOutput: readFixture("hosts_template.txt"),
		},
		{
			name:           "get hosts metrics with pagination",
			args:           []string{"--per-page", "1", "--page", "2"},
			metrics:        hostsMetrics,
			expectedOutput: readFixture("hosts_page.txt"),
		},
		{
			name:           "get all hosts metrics",
			args:           []string{"--per-page", "1", "-A"},
			metrics:        hostsMetrics,
			expectedOutput: readFixture("hosts.txt"),
		},
		{
			name:           "get empty hosts metrics",
			metrics:        "",
			expectedOutput: readFixture("hosts_empty.txt"),
		},
		{
			name:           "get hosts metrics in raw format",
			args:           []string{"--output", "raw"},
			metrics:        hostsMetrics,
			expectedOutput: hostsMetrics,
		},
		{
			name:        "get hosts metrics in unsupported format",
			args:        []string{"--output", "json"},
			noAPICall:   true,
			expectError: true,
		},
		{
			name:        "get hosts metrics in raw format with pagination",
			args:        []string{"--output", "raw", "--page", "2"},
			noAPICall:   true,
			expectError: true,
		},
		{
			name:        "get all hosts metrics in raw format",
			args:        []string{"--output", "raw", "-A"},
			noAPICall:   true,
			expectError: true,
		},
		{
			name:        "get hosts metrics with error",
			apiError:    true,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			metricsServiceHandler := mocks.NewMockMetricsService(mockCtrl)
			scClient := serverscom.NewClientWithEndpoint("", "")
			scClient.Metrics = metricsServiceHandler

			if !tc.noAPICall {
				var apiErr error
				if tc.apiError {
					apiErr = errors.New("some error")
				}
				metricsServiceHandler.EXPECT().
					ListHostsMetrics(gomock.Any()).
					Return(tc.metrics, apiErr)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			metricsCmd := NewCmd(testCmdContext)

			builder := testutils.NewTestCommandBuilder().
				WithCommand(metricsCmd).
				WithArgs(append([]string{"metrics", "hosts"}, tc.args...))

			cmd := builder.Build()
			var errOut bytes.Buffer
			cmd.SetErr(&errOut)

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
			g.Expect(builder.GetOutput()).To(BeEquivalentTo(tc.expectedOutput))
			g.Expect(errOut.String()).To(BeEquivalentTo(tc.expectedErrOut))
		})
	}
}

func TestRacksCmd(t *testing.T) {
	racksMetrics := readFixture("racks_input.txt")

	testCases := []metricsTestCase{
		{
			name:           "get racks metrics in default format",
			metrics:        racksMetrics,
			expectedOutput: readFixture("racks.txt"),
		},
		{
			name:           "get racks metrics in page view",
			args:           []string{"--page-view"},
			metrics:        racksMetrics,
			expectedOutput: readFixture("racks_page_view.txt"),
		},
		{
			name:           "get racks metrics with fields",
			args:           []string{"-f", "RackID", "-f", "AtsWatts", "-f", "AtsAmperes", "-f", "AtsCount"},
			metrics:        racksMetrics,
			expectedOutput: readFixture("racks_field.txt"),
		},
		{
			name:           "get racks metrics with pagination",
			args:           []string{"--per-page", "1"},
			metrics:        racksMetrics,
			expectedOutput: readFixture("racks_page.txt"),
		},
		{
			name:           "get all racks metrics",
			args:           []string{"--per-page", "1", "-A"},
			metrics:        racksMetrics,
			expectedOutput: readFixture("racks.txt"),
		},
		{
			name:           "get empty racks metrics",
			metrics:        "",
			expectedOutput: readFixture("racks_empty.txt"),
		},
		{
			name:           "get racks metrics in raw format",
			args:           []string{"--output", "raw"},
			metrics:        racksMetrics,
			expectedOutput: racksMetrics,
		},
		{
			name:        "get racks metrics in unsupported format",
			args:        []string{"--output", "yaml"},
			noAPICall:   true,
			expectError: true,
		},
		{
			name:        "get racks metrics with error",
			apiError:    true,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			metricsServiceHandler := mocks.NewMockMetricsService(mockCtrl)
			scClient := serverscom.NewClientWithEndpoint("", "")
			scClient.Metrics = metricsServiceHandler

			if !tc.noAPICall {
				var apiErr error
				if tc.apiError {
					apiErr = errors.New("some error")
				}
				metricsServiceHandler.EXPECT().
					ListRacksMetrics(gomock.Any()).
					Return(tc.metrics, apiErr)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			metricsCmd := NewCmd(testCmdContext)

			builder := testutils.NewTestCommandBuilder().
				WithCommand(metricsCmd).
				WithArgs(append([]string{"metrics", "racks"}, tc.args...))

			cmd := builder.Build()
			var errOut bytes.Buffer
			cmd.SetErr(&errOut)

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
			g.Expect(builder.GetOutput()).To(BeEquivalentTo(tc.expectedOutput))
			g.Expect(errOut.String()).To(BeEquivalentTo(tc.expectedErrOut))
		})
	}
}
