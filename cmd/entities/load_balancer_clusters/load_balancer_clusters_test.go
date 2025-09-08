package loadbalancerclusters

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	testId          = "testId"
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "lb-clusters")
	fixedTime       = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testLbCluster   = serverscom.LoadBalancerCluster{
		ID:         testId,
		Name:       "test-lb-cluster",
		LocationID: 1,
		Status:     "active",
		Created:    fixedTime,
		Updated:    fixedTime,
	}
)

func TestGetLoadBalancerClusterCmd(t *testing.T) {
	testCases := []struct {
		id             string
		name           string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get LB Cluster in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get LB Cluster in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get LB Cluster in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get LB Cluster with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbClustersServiceHandler := mocks.NewMockLoadBalancerClustersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancerClusters = lbClustersServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			lbClustersServiceHandler.EXPECT().
				GetLoadBalancerCluster(gomock.Any(), testId).
				Return(&testLbCluster, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbClusterCmd := NewCmd(testCmdContext)

			args := []string{"lb-clusters", "get", fmt.Sprint(tc.id)}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbClusterCmd).
				WithArgs(args)

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestListLoadBalancerClustersCmd(t *testing.T) {
	testLbCluster1 := testLbCluster
	testLbCluster2 := testLbCluster

	testLbCluster1.ID += "1"
	testLbCluster1.Name = "test-lb-cluster-1"

	testLbCluster2.ID += "2"
	testLbCluster2.Name = "test-lb-cluster-2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.LoadBalancerCluster])
	}{
		{
			name:           "list all LB Clusters",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancerCluster]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.LoadBalancerCluster{
						testLbCluster1,
						testLbCluster2,
					}, nil)
			},
		},
		{
			name:           "list LB Clusters",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancerCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.LoadBalancerCluster{
						testLbCluster1,
					}, nil)
			},
		},
		{
			name:           "list LB Clusters with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancerCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.LoadBalancerCluster{
						testLbCluster1,
						testLbCluster2,
					}, nil)
			},
		},
		{
			name:           "list LB Clusters with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancerCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.LoadBalancerCluster{
						testLbCluster1,
						testLbCluster2,
					}, nil)
			},
		},
		{
			name:        "list LB Clusters with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancerCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbClustersServiceHandler := mocks.NewMockLoadBalancerClustersService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.LoadBalancerCluster](mockCtrl)

	lbClustersServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancerClusters = lbClustersServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbClusterCmd := NewCmd(testCmdContext)

			args := []string{"lb-clusters", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbClusterCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}
