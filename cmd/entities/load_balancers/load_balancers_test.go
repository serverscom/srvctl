package loadbalancers

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
	fixtureBasePath      = filepath.Join("..", "..", "..", "testdata", "entities", "lb")
	skeletonTemplatePath = filepath.Join("..", "..", "..", "internal", "output", "skeletons", "templates", "lb")
	fixedTime            = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testId               = "testId"
	testLB               = serverscom.LoadBalancer{
		ID:                testId,
		Name:              "test-l4-lb",
		Type:              "l4",
		LocationID:        1,
		ExternalAddresses: []string{"127.0.0.1"},
		Status:            "active",
		Labels:            map[string]string{"foo": "bar"},
		Created:           fixedTime,
		Updated:           fixedTime,
	}
	testL4LB = serverscom.L4LoadBalancer{
		ID:                testId,
		Name:              "test-l4-lb",
		Type:              "l4",
		LocationID:        1,
		ExternalAddresses: []string{"127.0.0.1"},
		Status:            "active",
		Labels:            map[string]string{"foo": "bar"},
		Created:           fixedTime,
		Updated:           fixedTime,
	}
	testL7LB = serverscom.L7LoadBalancer{
		ID:                testId,
		Name:              "test-l7-lb",
		Type:              "l7",
		LocationID:        1,
		ExternalAddresses: []string{"127.0.0.1"},
		Domains:           []string{"test.com"},
		Status:            "active",
		Labels:            map[string]string{"foo": "bar"},
		Created:           fixedTime,
		Updated:           fixedTime,
	}
)

func TestAddL4LBCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockLoadBalancersService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create l4 lb",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l4.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create_l4.json")},
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					CreateL4LoadBalancer(gomock.Any(), gomock.AssignableToTypeOf(serverscom.L4LoadBalancerCreateInput{})).
					Return(&testL4LB, nil)
			},
		},
		{
			name:           "skeleton for l4 lb input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add_l4.json")),
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					CreateL4LoadBalancer(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:        "with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(lbServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l4", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestAddL7LBCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockLoadBalancersService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create l7 lb",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l7.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create_l7.json")},
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					CreateL7LoadBalancer(gomock.Any(), gomock.AssignableToTypeOf(serverscom.L7LoadBalancerCreateInput{})).
					Return(&testL7LB, nil)
			},
		},
		{
			name:           "skeleton for l7 lb input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add_l7.json")),
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					CreateL7LoadBalancer(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:        "with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(lbServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l7", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestGetL4LBCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get l4 lb in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l4.txt")),
		},
		{
			name:           "get l4 lb in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l4.json")),
		},
		{
			name:           "get l4 lb in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l4.yaml")),
		},
		{
			name:        "get l4 lb with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			lbServiceHandler.EXPECT().
				GetL4LoadBalancer(gomock.Any(), tc.id).
				Return(&testL4LB, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l4", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
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

func TestGetL7LBCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get l7 lb in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l7.txt")),
		},
		{
			name:           "get l7 lb in JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l7.json")),
		},
		{
			name:           "get l7 lb in YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_l7.yaml")),
		},
		{
			name:        "get l7 lb with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			lbServiceHandler.EXPECT().
				GetL7LoadBalancer(gomock.Any(), tc.id).
				Return(&testL7LB, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l7", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
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

func TestListLBCmd(t *testing.T) {
	lb1 := testLB
	lb1.ID += "1"
	lb2 := testLB
	lb2.ID += "2"
	lb2.Name = "test-l7-lb"
	lb2.Type = "l7"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.LoadBalancer])
	}{
		{
			name:           "list all load balancers",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancer]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.LoadBalancer{
						lb1,
						lb2,
					}, nil)
			},
		},
		{
			name:           "list l4 lb",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_l4.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.LoadBalancer{
						lb1,
					}, nil)
			},
		},
		{
			name:           "list l7 lb",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_l7.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.LoadBalancer{
						lb2,
					}, nil)
			},
		},
		{
			name:           "list lb with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.LoadBalancer{
						lb1,
						lb2,
					}, nil)
			},
		},
		{
			name:           "list lb with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.LoadBalancer{
						lb1,
						lb2,
					}, nil)
			},
		},
		{
			name:        "list lb with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.LoadBalancer]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.LoadBalancer](mockCtrl)

	lbServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
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

func TestUpdateL4LBCmd(t *testing.T) {
	updatedL4LB := testL4LB
	updatedL4LB.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockLoadBalancersService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update l4 lb",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_l4.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "update_input_l4.json")},
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					UpdateL4LoadBalancer(gomock.Any(), testId, serverscom.L4LoadBalancerUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&updatedL4LB, nil)
			},
		},
		{
			name:           "skeleton for update l4 lb input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "update_l4.json")),
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					UpdateL4LoadBalancer(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name: "update l4 lb with error",
			id:   testId,
			args: []string{"--input", filepath.Join(fixtureBasePath, "update_input_l4.json")},
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					UpdateL4LoadBalancer(gomock.Any(), testId, gomock.AssignableToTypeOf(serverscom.L4LoadBalancerUpdateInput{})).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(lbServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l4", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestUpdateL7LBCmd(t *testing.T) {
	updatedL7LB := testL7LB
	updatedL7LB.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockLoadBalancersService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update l7 lb",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_l7.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "update_input_l7.json")},
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					UpdateL7LoadBalancer(gomock.Any(), testId, serverscom.L7LoadBalancerUpdateInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&updatedL7LB, nil)
			},
		},
		{
			name:           "skeleton for update l7 lb input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "update_l7.json")),
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					UpdateL7LoadBalancer(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name: "update l7 lb with error",
			id:   testId,
			args: []string{"--input", filepath.Join(fixtureBasePath, "update_input_l7.json")},
			configureMock: func(mock *mocks.MockLoadBalancersService) {
				mock.EXPECT().
					UpdateL7LoadBalancer(gomock.Any(), testId, gomock.AssignableToTypeOf(serverscom.L7LoadBalancerUpdateInput{})).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(lbServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l7", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestDeleteL4LBCmd(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name: "delete l4 lb",
			id:   testId,
		},
		{
			name:        "delete l4 lb with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			lbServiceHandler.EXPECT().
				DeleteL4LoadBalancer(gomock.Any(), testId).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l4", "delete", tc.id}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
				WithArgs(args)

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestDeleteL7LBCmd(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name: "delete l7 lb",
			id:   testId,
		},
		{
			name:        "delete l7 lb with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	lbServiceHandler := mocks.NewMockLoadBalancersService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.LoadBalancers = lbServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			lbServiceHandler.EXPECT().
				DeleteL7LoadBalancer(gomock.Any(), testId).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			lbCmd := NewCmd(testCmdContext)

			args := []string{"lb", "l7", "delete", tc.id}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(lbCmd).
				WithArgs(args)

			cmd := builder.Build()

			err = cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
