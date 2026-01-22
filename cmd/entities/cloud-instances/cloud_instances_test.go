package cloudinstances

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
	fixtureBasePath     = filepath.Join("..", "..", "..", "testdata", "entities", "cloud-instances")
	fixedTime           = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testCloudInstanceID = "test-instance-id"
	testCloudInstance   = serverscom.CloudComputingInstance{
		ID:                 testCloudInstanceID,
		Name:               "test-instance",
		RegionID:           1,
		RegionCode:         "AMS1",
		OpenstackUUID:      "uuid-123",
		Status:             "active",
		FlavorID:           "flavor-1",
		FlavorName:         "m1.small",
		ImageID:            "image-1",
		ImageName:          stringPtr("Ubuntu 20.04"),
		PublicIPv4Address:  stringPtr("1.2.3.4"),
		PrivateIPv4Address: stringPtr("10.0.0.1"),
		LocalIPv4Address:   stringPtr("192.168.0.1"),
		PublicIPv6Address:  stringPtr("2001:db8::1"),
		GPNEnabled:         true,
		IPv6Enabled:        true,
		BackupCopies:       2,
		PublicPortBlocked:  false,
		Labels:             map[string]string{"env": "test"},
		Created:            fixedTime,
		Updated:            fixedTime,
	}
	testCloudInstance2 = serverscom.CloudComputingInstance{
		ID:                 "test-instance-id2",
		Name:               "test-instance2",
		RegionID:           1,
		RegionCode:         "AMS1",
		OpenstackUUID:      "uuid-123",
		Status:             "active",
		FlavorID:           "flavor-1",
		FlavorName:         "m1.small",
		ImageID:            "image-1",
		ImageName:          stringPtr("Ubuntu 20.04"),
		PublicIPv4Address:  stringPtr("1.2.3.4"),
		PrivateIPv4Address: stringPtr("10.0.0.1"),
		LocalIPv4Address:   stringPtr("192.168.0.1"),
		PublicIPv6Address:  stringPtr("2001:db8::1"),
		GPNEnabled:         true,
		IPv6Enabled:        true,
		BackupCopies:       2,
		PublicPortBlocked:  false,
		Labels:             map[string]string{"env": "test"},
		Created:            fixedTime,
		Updated:            fixedTime,
	}
)

func stringPtr(s string) *string {
	return &s
}

func TestListCloudInstancesCmd(t *testing.T) {
	testInstance1 := testCloudInstance
	testInstance2 := testCloudInstance2
	testInstance2.Name = "test-instance 2"
	testInstance2.Status = "inactive"
	testInstance2.RegionCode = "FRA1"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.CloudComputingInstance])
	}{
		{
			name:           "list all cloud instances",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingInstance]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.CloudComputingInstance{
						testInstance1,
						testInstance2,
					}, nil)
			},
		},
		{
			name:           "list cloud instances",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingInstance]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingInstance{
						testInstance1,
					}, nil)
			},
		},
		{
			name:           "list cloud instances with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingInstance]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingInstance{
						testInstance1,
						testInstance2,
					}, nil)
			},
		},
		{
			name:           "list cloud instances with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingInstance]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingInstance{
						testInstance1,
						testInstance2,
					}, nil)
			},
		},
		{
			name:        "list cloud instances with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingInstance]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	collectionHandler := mocks.NewMockCollection[serverscom.CloudComputingInstance](mockCtrl)
	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	cloudServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestGetCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		instanceID     string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get cloud instance in default format",
			instanceID:     testCloudInstanceID,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get cloud instance in JSON format",
			instanceID:     testCloudInstanceID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get cloud instance in YAML format",
			instanceID:     testCloudInstanceID,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			cloudServiceHandler.EXPECT().
				Get(gomock.Any(), tc.instanceID).
				Return(&testCloudInstance, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "get", "--instance-id", tc.instanceID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestAddCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudComputingInstancesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create cloud instance with input",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create.json")},
			configureMock: func(mock *mocks.MockCloudComputingInstancesService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.CloudComputingInstanceCreateInput{
						Name:     "test-instance",
						FlavorID: "flavor-1",
						ImageID:  "image-1",
						RegionID: 1,
						Labels:   map[string]string{"env": "test"},
					}).
					Return(&testCloudInstance, nil)
			},
		},
		{
			name:        "create cloud instance with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(cloudServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestUpdateCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		instanceID     string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudComputingInstancesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update cloud instance",
			instanceID:     testCloudInstanceID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args: []string{
				"--name", "updated-instance",
			},
			configureMock: func(mock *mocks.MockCloudComputingInstancesService) {
				mock.EXPECT().
					Update(gomock.Any(), testCloudInstanceID, gomock.AssignableToTypeOf(serverscom.CloudComputingInstanceUpdateInput{})).
					Return(&testCloudInstance, nil)
			},
		},
		{
			name:        "update cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
			configureMock: func(mock *mocks.MockCloudComputingInstancesService) {
				mock.EXPECT().
					Update(gomock.Any(), testCloudInstanceID, gomock.AssignableToTypeOf(serverscom.CloudComputingInstanceUpdateInput{})).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(cloudServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "update", tc.instanceID}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestDeleteCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name        string
		instanceID  string
		expectError bool
	}{
		{
			name:       "delete cloud instance",
			instanceID: testCloudInstanceID,
		},
		{
			name:        "delete cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			cloudServiceHandler.EXPECT().
				Delete(gomock.Any(), tc.instanceID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "delete", tc.instanceID}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestReinstallCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		instanceID     string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "reinstall cloud instance",
			instanceID:     testCloudInstanceID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:        "reinstall cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var instance *serverscom.CloudComputingInstance
			var err error
			if tc.expectError {
				err = errors.New("some error")
			} else {
				instance = &testCloudInstance
			}
			cloudServiceHandler.EXPECT().
				Reinstall(gomock.Any(), tc.instanceID, gomock.AssignableToTypeOf(serverscom.CloudComputingInstanceReinstallInput{})).
				Return(instance, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "reinstall", tc.instanceID, "--image-id", "image-1"}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestUpgradeCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		instanceID     string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "upgrade cloud instance",
			instanceID:     testCloudInstanceID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:        "upgrade cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var instance *serverscom.CloudComputingInstance
			var err error
			if tc.expectError {
				err = errors.New("some error")
			} else {
				instance = &testCloudInstance
			}
			cloudServiceHandler.EXPECT().
				Upgrade(gomock.Any(), tc.instanceID, gomock.AssignableToTypeOf(serverscom.CloudComputingInstanceUpgradeInput{})).
				Return(instance, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "upgrade", tc.instanceID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestRebootCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name        string
		instanceID  string
		expectError bool
	}{
		{
			name:       "reboot cloud instance",
			instanceID: testCloudInstanceID,
		},
		{
			name:        "reboot cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			cloudServiceHandler.EXPECT().
				Reboot(gomock.Any(), tc.instanceID).
				Return(nil, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "reboot", tc.instanceID}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestUpgradeRevertCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		instanceID     string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "upgrade revert cloud instance",
			instanceID:     testCloudInstanceID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:        "upgrade revert cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var instance *serverscom.CloudComputingInstance
			var err error
			if tc.expectError {
				err = errors.New("some error")
			} else {
				instance = &testCloudInstance
			}
			cloudServiceHandler.EXPECT().
				RevertUpgrade(gomock.Any(), tc.instanceID).
				Return(instance, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "upgrade-revert", tc.instanceID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestUpgradeApproveCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		instanceID     string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "upgrade approve cloud instance",
			instanceID:     testCloudInstanceID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:        "upgrade approve cloud instance with error",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var instance *serverscom.CloudComputingInstance
			var err error
			if tc.expectError {
				err = errors.New("some error")
			} else {
				instance = &testCloudInstance
			}
			cloudServiceHandler.EXPECT().
				ApproveUpgrade(gomock.Any(), tc.instanceID).
				Return(instance, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "upgrade-approve", tc.instanceID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestRescueModeCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name        string
		action      string
		instanceID  string
		expectError bool
	}{
		{
			name:       "enable rescue mode for cloud instance",
			action:     "activate",
			instanceID: testCloudInstanceID,
		},
		{
			name:        "enable rescue mode for cloud instance with error",
			action:      "activate",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
		{
			name:       "disable rescue mode for cloud instance",
			action:     "deactivate",
			instanceID: testCloudInstanceID,
		},
		{
			name:        "disable rescue mode for cloud instance with error",
			action:      "deactivate",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			switch tc.action {
			case "activate":
				cloudServiceHandler.EXPECT().
					Rescue(gomock.Any(), tc.instanceID).
					Return(&testCloudInstance, err)
			case "deactivate":
				cloudServiceHandler.EXPECT().
					Unrescue(gomock.Any(), tc.instanceID).
					Return(&testCloudInstance, err)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "rescue-mode", tc.instanceID, "--command", tc.action}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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

func TestPowerCloudInstancesCmd(t *testing.T) {
	testCases := []struct {
		name        string
		action      string
		instanceID  string
		expectError bool
	}{
		{
			name:       "power on cloud instance",
			action:     "on",
			instanceID: testCloudInstanceID,
		},
		{
			name:        "power on cloud instance with error",
			action:      "on",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
		{
			name:       "power off cloud instance",
			action:     "off",
			instanceID: testCloudInstanceID,
		},
		{
			name:        "power off cloud instance with error",
			action:      "off",
			instanceID:  testCloudInstanceID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudServiceHandler := mocks.NewMockCloudComputingInstancesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			switch tc.action {
			case "on":
				cloudServiceHandler.EXPECT().
					PowerOn(gomock.Any(), tc.instanceID).
					Return(&testCloudInstance, err)
			case "off":
				cloudServiceHandler.EXPECT().
					PowerOff(gomock.Any(), tc.instanceID).
					Return(&testCloudInstance, err)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudCmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "power", tc.instanceID, "--command", tc.action}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudCmd).
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
