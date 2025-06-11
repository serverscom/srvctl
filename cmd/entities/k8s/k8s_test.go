package k8s

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"fmt"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	testId                = "testId"
	testNodeId            = "testNodeId"
	fixtureBasePath       = filepath.Join("..", "..", "..", "testdata", "entities", "k8s")
	fixedTime             = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testKubernetesCluster = serverscom.KubernetesCluster{
		ID:         testId,
		Name:       "test-cluster",
		Status:     "active",
		LocationID: 1,
		Created:    fixedTime,
		Updated:    fixedTime,
	}
	testKubernetesClusterNode = serverscom.KubernetesClusterNode{
		ID:                 testNodeId,
		ClusterID:          testId,
		Number:             1,
		Hostname:           "test-node-1",
		Type:               "cloud",
		Role:               "master",
		Status:             "active",
		PrivateIPv4Address: "10.0.0.1",
		PublicIPv4Address:  "127.0.0.1",
		RefID:              "1",
		Configuration:      "SSD.50",
		Created:            fixedTime,
		Updated:            fixedTime,
	}
)

func TestGetKubernetesClusterCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get k8s cluster in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get k8s cluster in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get k8s cluster in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get k8s cluster with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	k8sServiceHandler := mocks.NewMockKubernetesClustersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.KubernetesClusters = k8sServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			k8sServiceHandler.EXPECT().
				Get(gomock.Any(), testId).
				Return(&testKubernetesCluster, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			k8sCmd := NewCmd(testCmdContext)

			args := []string{"k8s", "get", fmt.Sprint(tc.id)}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(k8sCmd).
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

func TestListKubernetesClustersCmd(t *testing.T) {
	testK8sCluster1 := testKubernetesCluster
	testK8sCluster2 := testKubernetesCluster
	testK8sCluster1.ID += "1"
	testK8sCluster2.Name = "test-cluster 2"
	testK8sCluster2.ID += "2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.KubernetesCluster])
	}{
		{
			name:           "list all k8s clusters",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesCluster]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.KubernetesCluster{
						testK8sCluster1,
						testK8sCluster2,
					}, nil)
			},
		},
		{
			name:           "list k8s clusters",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesCluster{
						testK8sCluster1,
					}, nil)
			},
		},
		{
			name:           "list k8s clusters with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesCluster{
						testK8sCluster1,
						testK8sCluster2,
					}, nil)
			},
		},
		{
			name:           "list k8s clusters with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesCluster{
						testK8sCluster1,
						testK8sCluster2,
					}, nil)
			},
		},
		{
			name:        "list k8s clusters with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesCluster]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	k8sServiceHandler := mocks.NewMockKubernetesClustersService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.KubernetesCluster](mockCtrl)

	k8sServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.KubernetesClusters = k8sServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			k8sCmd := NewCmd(testCmdContext)

			args := []string{"k8s", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(k8sCmd).
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

func TestUpdateKubernetesClusteCmd(t *testing.T) {
	newCluster := testKubernetesCluster
	newCluster.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockKubernetesClustersService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update k8s cluster",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockKubernetesClustersService) {
				mock.EXPECT().
					Update(gomock.Any(), testId, serverscom.KubernetesClusterUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newCluster, nil)
			},
		},
		{
			name: "update k8s cluster with error",
			id:   testId,
			configureMock: func(mock *mocks.MockKubernetesClustersService) {
				mock.EXPECT().
					Update(gomock.Any(), testId, serverscom.KubernetesClusterUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	k8sServiceHandler := mocks.NewMockKubernetesClustersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.KubernetesClusters = k8sServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(k8sServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			k8sCmd := NewCmd(testCmdContext)

			args := []string{"k8s", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(k8sCmd).
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

func TestGetKubernetesClusterNodeCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get k8s cluster node in default format",
			id:             testNodeId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node.txt")),
		},
		{
			name:           "get k8s cluster node in JSON format",
			id:             testNodeId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node.json")),
		},
		{
			name:           "get k8s cluster node in YAML format",
			id:             testNodeId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_node.yaml")),
		},
		{
			name:        "get k8s cluster node with error",
			id:          testNodeId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	k8sServiceHandler := mocks.NewMockKubernetesClustersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.KubernetesClusters = k8sServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			k8sServiceHandler.EXPECT().
				GetNode(gomock.Any(), testId, testNodeId).
				Return(&testKubernetesClusterNode, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			k8sCmd := NewCmd(testCmdContext)

			args := []string{"k8s", "get-node", fmt.Sprint(tc.id), "--cluster-id", testId}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(k8sCmd).
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

func TestListKubernetesClusterNodesCmd(t *testing.T) {
	testK8sClusterNode1 := testKubernetesClusterNode
	testK8sClusterNode2 := testKubernetesClusterNode
	testK8sClusterNode1.ID += "1"
	testK8sClusterNode2.Hostname = "test-node-2"
	testK8sClusterNode2.ID += "2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.KubernetesClusterNode])
	}{
		{
			name:           "list all k8s clusters",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all_nodes.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesClusterNode]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.KubernetesClusterNode{
						testK8sClusterNode1,
						testK8sClusterNode2,
					}, nil)
			},
		},
		{
			name:           "list k8s clusters",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_nodes.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesClusterNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesClusterNode{
						testK8sClusterNode1,
					}, nil)
			},
		},
		{
			name:           "list k8s clusters with template",
			args:           []string{"--template", "{{range .}}Hostname: {{.Hostname}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template_nodes.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesClusterNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesClusterNode{
						testK8sClusterNode1,
						testK8sClusterNode2,
					}, nil)
			},
		},
		{
			name:           "list k8s clusters with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview_nodes.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesClusterNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.KubernetesClusterNode{
						testK8sClusterNode1,
						testK8sClusterNode2,
					}, nil)
			},
		},
		{
			name:        "list k8s clusters with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.KubernetesClusterNode]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	k8sServiceHandler := mocks.NewMockKubernetesClustersService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.KubernetesClusterNode](mockCtrl)

	k8sServiceHandler.EXPECT().
		Nodes(testId).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.KubernetesClusters = k8sServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			k8sCmd := NewCmd(testCmdContext)

			args := []string{"k8s", "list-nodes", testId}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(k8sCmd).
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
