package cloudregions

import (
	"errors"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"
	"go.uber.org/mock/gomock"
)

var (
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "cloud-regions")
	testRegionID    = int64(1)
	testRegion      = serverscom.CloudComputingRegion{
		ID:   1,
		Name: "Test Region",
		Code: "TEST",
	}
	testRegion2 = serverscom.CloudComputingRegion{
		ID:   2,
		Name: "Test Region 2",
		Code: "TEST2",
	}
	testCredentials = serverscom.CloudComputingRegionCredentials{
		Password:   "test-password",
		TenantName: 12345,
		URL:        "https://test.example.com",
		Username:   67890,
	}
	testImage = serverscom.CloudComputingImage{
		ID:   "img-123",
		Name: "Test Image",
	}
	testImage2 = serverscom.CloudComputingImage{
		ID:   "img-456",
		Name: "Test Image 2",
	}
	testFlavor = serverscom.CloudComputingFlavor{
		ID:   "flavor-123",
		Name: "Test Flavor",
	}
	testFlavor2 = serverscom.CloudComputingFlavor{
		ID:   "flavor-456",
		Name: "Test Flavor 2",
	}
	testSnapshot = serverscom.CloudSnapshot{
		ID:        "snap-123",
		Name:      "Test Snapshot",
		ImageSize: 1024,
		MinDisk:   10,
		Status:    "available",
		IsBackup:  false,
		FileURL:   "https://test.example.com/snap-123",
	}
	testSnapshot2 = serverscom.CloudSnapshot{
		ID:        "snap-456",
		Name:      "Test Snapshot 2",
		ImageSize: 1024,
		MinDisk:   10,
		Status:    "available",
		IsBackup:  false,
		FileURL:   "https://test.example.com/snap-123",
	}
)

func TestListCloudRegionsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.CloudComputingRegion])
	}{
		{
			name:           "list all cloud regions",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingRegion]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.CloudComputingRegion{
						testRegion,
						testRegion2,
					}, nil)
			},
		},
		{
			name:           "list cloud regions",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingRegion]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingRegion{
						testRegion,
					}, nil)
			},
		},
		{
			name:           "list cloud regions with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingRegion]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingRegion{
						testRegion,
						testRegion2,
					}, nil)
			},
		},
		{
			name:           "list cloud regions with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingRegion]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingRegion{
						testRegion,
						testRegion2,
					}, nil)
			},
		},
		{
			name:        "list cloud regions with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingRegion]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudRegionsServiceHandler := mocks.NewMockCloudComputingRegionsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.CloudComputingRegion](mockCtrl)

	cloudRegionsServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingRegions = cloudRegionsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudRegionsCmd := NewCmd(testCmdContext)

			args := []string{"cloud-regions", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudRegionsCmd).
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

func TestGetCredentialsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		regionID       string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get credentials in default format",
			regionID:       "1",
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_credentials.txt")),
		},
		{
			name:           "get credentials in JSON format",
			regionID:       "1",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_credentials.json")),
		},
		{
			name:           "get credentials in YAML format",
			regionID:       "1",
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_credentials.yaml")),
		},
		{
			name:        "get credentials with error",
			regionID:    "1",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudRegionsServiceHandler := mocks.NewMockCloudComputingRegionsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingRegions = cloudRegionsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			cloudRegionsServiceHandler.EXPECT().
				Credentials(gomock.Any(), testRegionID).
				Return(&testCredentials, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudRegionsCmd := NewCmd(testCmdContext)

			args := []string{"cloud-regions", "get-credentials", tc.regionID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudRegionsCmd).
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

func TestListImagesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.CloudComputingImage])
	}{
		{
			name:           "list all images",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_images_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingImage]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.CloudComputingImage{
						testImage,
						testImage2,
					}, nil)
			},
		},
		{
			name:           "list images",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_images.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingImage]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingImage{
						testImage,
					}, nil)
			},
		},
		{
			name:           "list images with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_images_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingImage]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingImage{
						testImage,
						testImage2,
					}, nil)
			},
		},
		{
			name:           "list images with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_images_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingImage]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingImage{
						testImage,
						testImage2,
					}, nil)
			},
		},
		{
			name:        "list images with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingImage]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudRegionsServiceHandler := mocks.NewMockCloudComputingRegionsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.CloudComputingImage](mockCtrl)

	cloudRegionsServiceHandler.EXPECT().
		Images(testRegionID).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingRegions = cloudRegionsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudRegionsCmd := NewCmd(testCmdContext)

			args := []string{"cloud-regions", "list-images", "1"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudRegionsCmd).
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

func TestListFlavorsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.CloudComputingFlavor])
	}{
		{
			name:           "list all flavors",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_flavors_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingFlavor]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.CloudComputingFlavor{
						testFlavor,
						testFlavor2,
					}, nil)
			},
		},
		{
			name:           "list flavors",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_flavors.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingFlavor{
						testFlavor,
					}, nil)
			},
		},
		{
			name:           "list flavors with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_flavors_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingFlavor{
						testFlavor,
						testFlavor2,
					}, nil)
			},
		},
		{
			name:           "list flavors with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_flavors_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudComputingFlavor{
						testFlavor,
						testFlavor2,
					}, nil)
			},
		},
		{
			name:        "list flavors with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudComputingFlavor]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudRegionsServiceHandler := mocks.NewMockCloudComputingRegionsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.CloudComputingFlavor](mockCtrl)

	cloudRegionsServiceHandler.EXPECT().
		Flavors(testRegionID).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingRegions = cloudRegionsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudRegionsCmd := NewCmd(testCmdContext)

			args := []string{"cloud-regions", "list-flavors", "1"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudRegionsCmd).
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

func TestListSnapshotsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.CloudSnapshot])
	}{
		{
			name:           "list all snapshots",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_snapshots_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudSnapshot]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.CloudSnapshot{
						testSnapshot,
						testSnapshot2,
					}, nil)
			},
		},
		{
			name:           "list snapshots",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_snapshots.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudSnapshot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudSnapshot{
						testSnapshot,
					}, nil)
			},
		},
		{
			name:           "list snapshots with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_snapshots_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudSnapshot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudSnapshot{
						testSnapshot,
						testSnapshot2,
					}, nil)
			},
		},
		{
			name:           "list snapshots with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_snapshots_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudSnapshot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudSnapshot{
						testSnapshot,
						testSnapshot2,
					}, nil)
			},
		},
		{
			name:        "list snapshots with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudSnapshot]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudRegionsServiceHandler := mocks.NewMockCloudComputingRegionsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.CloudSnapshot](mockCtrl)

	cloudRegionsServiceHandler.EXPECT().
		Snapshots(testRegionID).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingRegions = cloudRegionsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudRegionsCmd := NewCmd(testCmdContext)

			args := []string{"cloud-regions", "list-snapshots", "1"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudRegionsCmd).
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

func TestAddSnapshotCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudComputingRegionsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "add snapshot",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "add_snapshot.json")),
			args: []string{
				"--name", "Test Snapshot",
				"--instance-id", "instance-123",
			},
			configureMock: func(mock *mocks.MockCloudComputingRegionsService) {
				mock.EXPECT().
					CreateSnapshot(gomock.Any(), testRegionID, serverscom.CloudSnapshotCreateInput{
						Name:       "Test Snapshot",
						InstanceID: "instance-123",
					}).
					Return(&testSnapshot, nil)
			},
		},
		{
			name:        "add snapshot with error",
			expectError: true,
			args: []string{
				"--name", "Test Snapshot",
				"--instance-id", "instance-123",
			},
			configureMock: func(mock *mocks.MockCloudComputingRegionsService) {
				mock.EXPECT().
					CreateSnapshot(gomock.Any(), testRegionID, serverscom.CloudSnapshotCreateInput{
						Name:       "Test Snapshot",
						InstanceID: "instance-123",
					}).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudRegionsServiceHandler := mocks.NewMockCloudComputingRegionsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingRegions = cloudRegionsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(cloudRegionsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudRegionsCmd := NewCmd(testCmdContext)

			args := []string{"cloud-regions", "add-snapshot", "1"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudRegionsCmd).
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

func TestDeleteSnapshotCmd(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		configureMock func(*mocks.MockCloudComputingRegionsService)
		expectError   bool
	}{
		{
			name: "delete snapshot",
			args: []string{
				"--snapshot-id", "snap-123",
			},
			configureMock: func(mock *mocks.MockCloudComputingRegionsService) {
				mock.EXPECT().
					DeleteSnapshot(gomock.Any(), testRegionID, "snap-123").
					Return(nil)
			},
		},
		{
			name: "delete snapshot with error",
			args: []string{
				"--snapshot-id", "snap-123",
			},
			expectError: true,
			configureMock: func(mock *mocks.MockCloudComputingRegionsService) {
				mock.EXPECT().
					DeleteSnapshot(gomock.Any(), testRegionID, "snap-123").
					Return(errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudRegionsServiceHandler := mocks.NewMockCloudComputingRegionsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingRegions = cloudRegionsServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(cloudRegionsServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cloudRegionsCmd := NewCmd(testCmdContext)

			args := []string{"cloud-regions", "delete-snapshot", "1"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(cloudRegionsCmd).
				WithArgs(args)

			cmd := builder.Build()

			err := cmd.Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
