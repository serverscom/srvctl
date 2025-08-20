package networkpools

import (
	"errors"
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
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "network-pools")
	testID          = "testId"
	testSubnetID    = "testSubnetId"
	fixedTime       = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	networkTitle = "testNetworkPool"
	subnetTitle  = "testSubnet"
	cidr         = "192.168.0.0/24"

	testNetworkPool = serverscom.NetworkPool{
		ID:          testID,
		Title:       &networkTitle,
		CIDR:        cidr,
		Type:        "private",
		LocationIDs: []int{1, 2},
		Labels: map[string]string{
			"environment": "testing",
		},
		Created: fixedTime,
		Updated: fixedTime,
	}

	testSubnet = serverscom.Subnetwork{
		ID:            testSubnetID,
		NetworkPoolID: testID,
		Title:         &subnetTitle,
		CIDR:          cidr,
		Attached:      false,
		InterfaceType: "private",
		Created:       fixedTime,
		Updated:       fixedTime,
	}
)

func TestGetNetworkPoolCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get network pool default format",
			id:             testID,
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get network pool JSON",
			id:             testID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get network pool YAML",
			id:             testID,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get network pool with error",
			id:          testID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			var err error
			if tc.expectError {
				err = errors.New("some error")
			}

			poolService.EXPECT().
				Get(gomock.Any(), tc.id).
				Return(&testNetworkPool, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"network-pools", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)
			c := builder.Build()
			err = c.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestGetSubnetCmd(t *testing.T) {
	testCases := []struct {
		name           string
		networkID      string
		subnetID       string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get subnet default",
			networkID:      testID,
			subnetID:       testSubnetID,
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_subnet.txt")),
		},
		{
			name:           "get subnet json",
			networkID:      testID,
			subnetID:       testSubnetID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_subnet.json")),
		},
		{
			name:           "get subnet YAML",
			networkID:      testID,
			subnetID:       testSubnetID,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_subnet.yaml")),
		},
		{
			name:        "get subnet with error",
			networkID:   testID,
			subnetID:    testSubnetID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			poolService.EXPECT().
				GetSubnetwork(gomock.Any(), tc.networkID, tc.subnetID).
				Return(&testSubnet, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"network-pools", "get-subnet", tc.networkID, "--network-id", tc.subnetID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)
			c := builder.Build()
			err = c.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestListNetworkPoolsCmd(t *testing.T) {
	title := "testNetworkPool2"
	p1 := testNetworkPool
	p2 := testNetworkPool
	p2.ID += "2"
	p2.Title = &title

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.NetworkPool])
	}{
		{
			name:           "list all network pools",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.NetworkPool]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.NetworkPool{p1, p2}, nil)
			},
		},
		{
			name:           "list network pools",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.NetworkPool]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.NetworkPool{p1}, nil)
			},
		},
		{
			name:           "list network pools with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Title: {{.Title}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.NetworkPool]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.NetworkPool{p1, p2}, nil)
			},
		},
		{
			name:           "list network pools with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.NetworkPool]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.NetworkPool{p1, p2}, nil)
			},
		},
		{
			name:        "list network pools error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.NetworkPool]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.NetworkPool](mockCtrl)

	poolService.EXPECT().Collection().Return(collection).AnyTimes()
	collection.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collection).AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(collection)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := append([]string{"network-pools", "list"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestListSubnetsCmd(t *testing.T) {
	title := "testSubnet2"
	s1 := testSubnet
	s2 := testSubnet
	s2.ID += "2"
	s2.Title = &title

	testCases := []struct {
		name           string
		networkID      string
		output         string
		expectedOutput []byte
		args           []string
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Subnetwork])
	}{
		{
			name:           "list all subnets",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_subnets_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Subnetwork]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Subnetwork{s1, s2}, nil)
			},
		},
		{
			name:           "list subnets",
			networkID:      testID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_subnets.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Subnetwork]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Subnetwork{s1}, nil)
			},
		},
		{
			name:           "list subnets with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Title: {{.Title}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_subnets_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Subnetwork]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Subnetwork{s1, s2}, nil)
			},
		},
		{
			name:           "list subnets with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_subnets_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Subnetwork]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Subnetwork{s1, s2}, nil)
			},
		},
		{
			name:        "list subnets error",
			networkID:   testID,
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Subnetwork]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.Subnetwork](mockCtrl)

	poolService.EXPECT().Subnetworks(gomock.Any()).Return(collection).AnyTimes()
	collection.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collection).AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(collection)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := append([]string{"network-pools", "list-subnets", tc.networkID}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)
			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestCreateSubnetCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		output         string
		expectedOutput []byte
		configureMock  func(*mocks.MockNetworkPoolsService)
		expectError    bool
	}{
		{
			name:           "create subnet",
			args:           []string{testID, "--title", subnetTitle, "--cidr", cidr},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_subnet.json")),
			configureMock: func(mock *mocks.MockNetworkPoolsService) {
				mock.EXPECT().
					CreateSubnetwork(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.SubnetworkCreateInput{})).
					Return(&testSubnet, nil)
			},
		},
		{
			name:        "create subnet error",
			args:        []string{testID, "--title", subnetTitle, "--cidr", cidr},
			expectError: true,
			configureMock: func(mock *mocks.MockNetworkPoolsService) {
				mock.EXPECT().
					CreateSubnetwork(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.SubnetworkCreateInput{})).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(poolService)
			}
			ctx := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(ctx)

			args := append([]string{"network-pools", "add-subnet"}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestUpdateNetworkPoolCmd(t *testing.T) {
	updatedPool := testNetworkPool
	title := "updated"
	updatedPool.Title = &title
	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockNetworkPoolsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update network pool",
			id:             testID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update.json")),
			args:           []string{"--title", "updated"},
			configureMock: func(mock *mocks.MockNetworkPoolsService) {
				mock.EXPECT().
					Update(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.NetworkPoolInput{})).
					Return(&updatedPool, nil)
			},
		},
		{
			name: "update network pool error",
			id:   testID,
			args: []string{"--title", "updated"},
			configureMock: func(mock *mocks.MockNetworkPoolsService) {
				mock.EXPECT().
					Update(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.NetworkPoolInput{})).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(poolService)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := append([]string{"network-pools", "update", tc.id}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)
			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestUpdateSubnetCmd(t *testing.T) {
	updatedSubnet := testSubnet
	title := "updated"
	updatedSubnet.Title = &title
	testCases := []struct {
		name           string
		networkID      string
		subnetID       string
		args           []string
		output         string
		configureMock  func(*mocks.MockNetworkPoolsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update subnet",
			networkID:      testID,
			subnetID:       testSubnetID,
			output:         "json",
			args:           []string{"--network-id", testSubnetID, "--title", subnetTitle},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_subnet.json")),
			configureMock: func(mock *mocks.MockNetworkPoolsService) {
				mock.EXPECT().
					UpdateSubnetwork(gomock.Any(), testID, testSubnetID, gomock.AssignableToTypeOf(serverscom.SubnetworkUpdateInput{})).
					Return(&updatedSubnet, nil)
			},
		},
		{
			name:      "update subnet error",
			networkID: testID,
			subnetID:  testSubnetID,
			args:      []string{"--network-id", testSubnetID, "--title", subnetTitle},
			configureMock: func(mock *mocks.MockNetworkPoolsService) {
				mock.EXPECT().
					UpdateSubnetwork(gomock.Any(), testID, testSubnetID, gomock.AssignableToTypeOf(serverscom.SubnetworkUpdateInput{})).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(poolService)
			}
			ctx := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(ctx)

			args := append([]string{"network-pools", "update-subnet", tc.networkID}, tc.args...)
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)

			c := builder.Build()
			err := c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestDeleteSubnetCmd(t *testing.T) {
	testCases := []struct {
		name        string
		networkID   string
		subnetID    string
		expectError bool
	}{
		{
			name:      "delete subnet",
			networkID: testID,
			subnetID:  testSubnetID,
		},
		{
			name:        "delete subnet error",
			networkID:   testID,
			subnetID:    testSubnetID,
			expectError: true,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	poolService := mocks.NewMockNetworkPoolsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.NetworkPools = poolService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			poolService.EXPECT().
				DeleteSubnetwork(gomock.Any(), tc.networkID, tc.subnetID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"network-pools", "delete", tc.networkID, "--network-id", tc.subnetID}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(cmd).
				WithArgs(args)
			c := builder.Build()
			err = c.Execute()
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
