package cloudvolumes

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
	testVolumeID    = "vol-12345"
	testInstanceID  = "instance-123"
	fixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "cloud-volumes")
	fixedTime       = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testVolume      = serverscom.CloudBlockStorageVolume{
		ID:          testVolumeID,
		Name:        "test-volume",
		RegionID:    1,
		Size:        100,
		Description: testutils.PtrString("Test volume"),
		Labels:      map[string]string{"foo": "bar"},
		Created:     &fixedTime,
	}
)

func TestGetCloudVolumesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get volume text",
			output:         "text",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
			args:           []string{testVolumeID},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().Get(gomock.Any(), testVolumeID).Return(&testVolume, nil)
			},
		},
		{
			name:           "get volume json",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args:           []string{testVolumeID, "--output", "json"},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().Get(gomock.Any(), testVolumeID).Return(&testVolume, nil)
			},
		},
		{
			name:           "get volume yaml",
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
			args:           []string{testVolumeID, "--output", "yaml"},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().Get(gomock.Any(), testVolumeID).Return(&testVolume, nil)
			},
		},
		{
			name:        "get volume error",
			expectError: true,
			args:        []string{testVolumeID},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().Get(gomock.Any(), testVolumeID).Return(nil, errors.New("not found"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	volumesServiceHandler := mocks.NewMockCloudBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageVolumes = volumesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(volumesServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			volumesCmd := NewCmd(testCmdContext)

			args := []string{"cloud-volumes", "get"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().WithCommand(volumesCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestListCloudVolumesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageVolumesService, *mocks.MockCollection[serverscom.CloudBlockStorageVolume])
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "list all",
			args:           []string{"-A", "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(serviceMock *mocks.MockCloudBlockStorageVolumesService, collectionMock *mocks.MockCollection[serverscom.CloudBlockStorageVolume]) {
				serviceMock.EXPECT().Collection().Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().Collect(gomock.Any()).Return([]serverscom.CloudBlockStorageVolume{testVolume}, nil)
			},
		},
		{
			name:           "list",
			args:           []string{"--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(serviceMock *mocks.MockCloudBlockStorageVolumesService, collectionMock *mocks.MockCollection[serverscom.CloudBlockStorageVolume]) {
				serviceMock.EXPECT().Collection().Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().List(gomock.Any()).Return([]serverscom.CloudBlockStorageVolume{testVolume}, nil)
			},
		},
		{
			name:           "template",
			args:           []string{"--template", "{{range .}}{{.ID}}{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(serviceMock *mocks.MockCloudBlockStorageVolumesService, collectionMock *mocks.MockCollection[serverscom.CloudBlockStorageVolume]) {
				serviceMock.EXPECT().Collection().Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().List(gomock.Any()).Return([]serverscom.CloudBlockStorageVolume{testVolume}, nil)
			},
		},
		{
			name:           "page-view",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(serviceMock *mocks.MockCloudBlockStorageVolumesService, collectionMock *mocks.MockCollection[serverscom.CloudBlockStorageVolume]) {
				serviceMock.EXPECT().Collection().Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().List(gomock.Any()).Return([]serverscom.CloudBlockStorageVolume{testVolume}, nil)
			},
		},
		{
			name:        "error",
			expectError: true,
			configureMock: func(serviceMock *mocks.MockCloudBlockStorageVolumesService, collectionMock *mocks.MockCollection[serverscom.CloudBlockStorageVolume]) {
				serviceMock.EXPECT().Collection().Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collectionMock).AnyTimes()
				collectionMock.EXPECT().List(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	volumesServiceHandler := mocks.NewMockCloudBlockStorageVolumesService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.CloudBlockStorageVolume](mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageVolumes = volumesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(volumesServiceHandler, collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			volumesCmd := NewCmd(testCmdContext)

			args := []string{"cloud-volumes", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().WithCommand(volumesCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestAddCloudVolumesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create volume with input",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create.json"), "--output", "json"},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.CloudBlockStorageVolumeCreateInput{
						Name:        "test-volume",
						RegionID:    1,
						Size:        100,
						Description: "Test volume",
						Labels:      map[string]string{"foo": "bar"},
					}).
					Return(&testVolume, nil)
			},
		},
		{
			name:           "create volume",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args: []string{
				"--name", "test-volume",
				"--region-id", "1",
				"--size", "100",
				"--description", "Test volume",
				"--label", "foo=bar",
				"--output", "json",
			},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.CloudBlockStorageVolumeCreateInput{
						Name:        "test-volume",
						RegionID:    1,
						Size:        100,
						Description: "Test volume",
						Labels:      map[string]string{"foo": "bar"},
					}).
					Return(&testVolume, nil)
			},
		},
		{
			name:        "create volume with error",
			expectError: true,
			args:        []string{"--name", "test-volume", "--region-id", "1"},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("create error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	volumesServiceHandler := mocks.NewMockCloudBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageVolumes = volumesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(volumesServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			volumesCmd := NewCmd(testCmdContext)

			args := []string{"cloud-volumes", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().WithCommand(volumesCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestUpdateCloudVolumesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update volume",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args: []string{
				testVolumeID,
				"--name", "updated-volume",
				"--description", "Updated volume",
				"--label", "new=label",
				"--output", "json",
			},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Update(gomock.Any(), testVolumeID, serverscom.CloudBlockStorageVolumeUpdateInput{
						Name:        "updated-volume",
						Description: "Updated volume",
						Labels:      map[string]string{"new": "label"},
					}).
					Return(&testVolume, nil)
			},
		},
		{
			name:        "update volume with error",
			expectError: true,
			args:        []string{testVolumeID, "--name", "updated-volume"},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Update(gomock.Any(), testVolumeID, gomock.Any()).
					Return(nil, errors.New("update error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	volumesServiceHandler := mocks.NewMockCloudBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageVolumes = volumesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(volumesServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			volumesCmd := NewCmd(testCmdContext)

			args := []string{"cloud-volumes", "update"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().WithCommand(volumesCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestDeleteCloudVolumesCmd(t *testing.T) {
	testCases := []struct {
		name        string
		volumeID    string
		expectError bool
	}{
		{
			name:     "delete volume",
			volumeID: testVolumeID,
		},
		{
			name:        "delete volume with error",
			volumeID:    testVolumeID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	volumesServiceHandler := mocks.NewMockCloudBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageVolumes = volumesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("delete error")
			}
			volumesServiceHandler.EXPECT().
				Delete(gomock.Any(), tc.volumeID).
				Return(nil, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			volumesCmd := NewCmd(testCmdContext)

			args := []string{"cloud-volumes", "delete", tc.volumeID}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(volumesCmd).
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

func TestAttachCloudVolumesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "attach volume",
			args:           []string{testVolumeID, "--instance-id", testInstanceID, "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Attach(gomock.Any(), testVolumeID, serverscom.CloudBlockStorageVolumeAttachInput{
						InstanceID: testInstanceID,
					}).
					Return(&testVolume, nil)
			},
		},
		{
			name:        "attach volume with error",
			expectError: true,
			args:        []string{testVolumeID, "--instance-id", testInstanceID},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Attach(gomock.Any(), testVolumeID, gomock.Any()).
					Return(nil, errors.New("attach error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	volumesServiceHandler := mocks.NewMockCloudBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageVolumes = volumesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(volumesServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			volumesCmd := NewCmd(testCmdContext)

			args := []string{"cloud-volumes", "volume-attach"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().WithCommand(volumesCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}

func TestDetachCloudVolumesCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "detach volume",
			args:           []string{testVolumeID, "--instance-id", testInstanceID, "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Detach(gomock.Any(), testVolumeID, serverscom.CloudBlockStorageVolumeDetachInput{
						InstanceID: testInstanceID,
					}).
					Return(&testVolume, nil)
			},
		},
		{
			name:        "detach volume with error",
			expectError: true,
			args:        []string{testVolumeID, "--instance-id", testInstanceID},
			configureMock: func(mock *mocks.MockCloudBlockStorageVolumesService) {
				mock.EXPECT().
					Detach(gomock.Any(), testVolumeID, gomock.Any()).
					Return(nil, errors.New("detach error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	volumesServiceHandler := mocks.NewMockCloudBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageVolumes = volumesServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(volumesServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			volumesCmd := NewCmd(testCmdContext)

			args := []string{"cloud-volumes", "volume-detach"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}

			builder := testutils.NewTestCommandBuilder().WithCommand(volumesCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(BeEquivalentTo(string(tc.expectedOutput)))
			}
		})
	}
}
