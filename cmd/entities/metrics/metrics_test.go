package metrics

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestHostsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		metrics        string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "get hosts metrics",
			metrics:        "# HELP hosts_total\nhosts_total 42\n",
			expectedOutput: "# HELP hosts_total\nhosts_total 42\n",
		},
		{
			name:        "get hosts metrics with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricsServiceHandler := mocks.NewMockMetricsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Metrics = metricsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			metricsServiceHandler.EXPECT().
				ListHostsMetrics(gomock.Any()).
				Return(tc.metrics, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			metricsCmd := NewCmd(testCmdContext)

			builder := testutils.NewTestCommandBuilder().
				WithCommand(metricsCmd).
				WithArgs([]string{"metrics", "hosts"})

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(Equal(tc.expectedOutput))
			}
		})
	}
}

func TestRacksCmd(t *testing.T) {
	testCases := []struct {
		name           string
		metrics        string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "get racks metrics",
			metrics:        "# HELP racks_total\nracks_total 7\n",
			expectedOutput: "# HELP racks_total\nracks_total 7\n",
		},
		{
			name:        "get racks metrics with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricsServiceHandler := mocks.NewMockMetricsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.Metrics = metricsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			metricsServiceHandler.EXPECT().
				ListRacksMetrics(gomock.Any()).
				Return(tc.metrics, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			metricsCmd := NewCmd(testCmdContext)

			builder := testutils.NewTestCommandBuilder().
				WithCommand(metricsCmd).
				WithArgs([]string{"metrics", "racks"})

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(Equal(tc.expectedOutput))
			}
		})
	}
}
