package rbsvolumes

import (
	"errors"
	_ "fmt"
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
	fixtureBasePath      = filepath.Join("..", "..", "..", "testdata", "entities", "rbs-volumes")
	skeletonTemplatePath = filepath.Join("..", "..", "..", "internal", "output", "skeletons", "templates", "rbs-volumes")
	fixedTime            = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testRBSVolumeID      = "testID"
	testRBSVolume        = serverscom.RemoteBlockStorageVolume{
		ID:        testRBSVolumeID,
		Name:      "rbs-volume",
		Labels:    map[string]string{"foo": "bar"},
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}
	testRBSCredentials = serverscom.RemoteBlockStorageVolumeCredentials{
		VolumeID: testRBSVolumeID,
		Username: "user",
		Password: "secret",
	}
)

func TestListRBSVolumesCmd(t *testing.T) {
	testRBSVolume1 := testRBSVolume
	testRBSVolume2 := testRBSVolume
	testRBSVolume1.ID += "1"
	testRBSVolume2.ID += "2"
	testRBSVolume1.Name += "1"
	testRBSVolume2.Name += "2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.RemoteBlockStorageVolume])
	}{
		{
			name:           "list all rbs volumes",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.RemoteBlockStorageVolume]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.RemoteBlockStorageVolume{
						testRBSVolume1,
						testRBSVolume2,
					}, nil)
			},
		},
		{
			name:           "list rbs volumes",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.RemoteBlockStorageVolume]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.RemoteBlockStorageVolume{
						testRBSVolume1,
					}, nil)
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rbsVolumeServiceHandler := mocks.NewMockRemoteBlockStorageVolumesService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.RemoteBlockStorageVolume](mockCtrl)

	rbsVolumeServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.RemoteBlockStorageVolumes = rbsVolumeServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rbsVolumeCmd := NewCmd(testCmdContext)

			args := []string{"rbs", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(rbsVolumeCmd).
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

func TestGetRBSVolumeCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockRemoteBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get rbs volume json",
			args:           []string{testRBSVolumeID, "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().Get(gomock.Any(), testRBSVolumeID).Return(&testRBSVolume, nil)
			},
		},
		{
			name:        "get rbs volume error",
			args:        []string{testRBSVolumeID},
			expectError: true,
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().Get(gomock.Any(), testRBSVolumeID).Return(nil, errors.New("not found"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rbsVolumeServiceHandler := mocks.NewMockRemoteBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.RemoteBlockStorageVolumes = rbsVolumeServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(rbsVolumeServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rbsVolumeCmd := NewCmd(testCmdContext)

			args := []string{"rbs", "get"}
			args = append(args, tc.args...)

			builder := testutils.NewTestCommandBuilder().WithCommand(rbsVolumeCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestAddRBSVolumeCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockRemoteBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "add rbs volume with input",
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create.json"), "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.RemoteBlockStorageVolumeCreateInput{
						Name:       "rbs-volume",
						Size:       100,
						LocationID: 1,
						FlavorID:   2,
						Labels:     map[string]string{"foo": "bar"},
					}).
					Return(&testRBSVolume, nil)
			},
		},
		{
			name: "add rbs volume with params",
			args: []string{
				"--name", "rbs-volume",
				"--size", "100",
				"--flavor-id", "2",
				"--location-id", "1",
				"--label", "foo=bar",
				"--output", "json",
			},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.RemoteBlockStorageVolumeCreateInput{
						Name:       "rbs-volume",
						Size:       100,
						LocationID: 1,
						FlavorID:   2,
						Labels:     map[string]string{"foo": "bar"},
					}).
					Return(&testRBSVolume, nil)
			},
		},
		{
			name:           "skeleton for rbs volume input",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add.json")),
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:        "add rbs volume with error",
			args:        []string{"--name", "rbs-volume", "--size", "100", "--flavor-id", "2"},
			expectError: true,
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("create error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rbsVolumeServiceHandler := mocks.NewMockRemoteBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.RemoteBlockStorageVolumes = rbsVolumeServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(rbsVolumeServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rbsVolumeCmd := NewCmd(testCmdContext)

			args := []string{"rbs", "add"}
			args = append(args, tc.args...)

			builder := testutils.NewTestCommandBuilder().WithCommand(rbsVolumeCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestUpdateRBSVolumeCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockRemoteBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name: "update rbs volume",
			args: []string{
				testRBSVolumeID,
				"--name", "updated-volume",
				"--size", "200",
				"--label", "new=label",
				"--output", "json",
			},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().
					Update(gomock.Any(), testRBSVolumeID, serverscom.RemoteBlockStorageVolumeUpdateInput{
						Name:   "updated-volume",
						Size:   200,
						Labels: map[string]string{"new": "label"},
					}).
					Return(&testRBSVolume, nil)
			},
		},
		{
			name:        "update rbs volume with error",
			args:        []string{testRBSVolumeID, "--name", "updated-volume"},
			expectError: true,
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().
					Update(gomock.Any(), testRBSVolumeID, gomock.Any()).
					Return(nil, errors.New("update error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rbsVolumeServiceHandler := mocks.NewMockRemoteBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.RemoteBlockStorageVolumes = rbsVolumeServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(rbsVolumeServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rbsVolumeCmd := NewCmd(testCmdContext)

			args := []string{"rbs", "update"}
			args = append(args, tc.args...)

			builder := testutils.NewTestCommandBuilder().WithCommand(rbsVolumeCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestDeleteRBSVolumeCmd(t *testing.T) {
	testCases := []struct {
		name        string
		volumeID    string
		expectError bool
	}{
		{
			name:     "delete rbs volume",
			volumeID: testRBSVolumeID,
		},
		{
			name:        "delete rbs volume with error",
			volumeID:    testRBSVolumeID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rbsVolumeServiceHandler := mocks.NewMockRemoteBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.RemoteBlockStorageVolumes = rbsVolumeServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var mockErr error
			if tc.expectError {
				mockErr = errors.New("delete error")
			}
			rbsVolumeServiceHandler.EXPECT().
				Delete(gomock.Any(), tc.volumeID).
				Return(mockErr)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rbsVolumeCmd := NewCmd(testCmdContext)

			args := []string{"rbs", "delete", tc.volumeID}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(rbsVolumeCmd).
				WithArgs(args)

			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestGetRBSVolumeCredentialsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockRemoteBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get credentials json",
			args:           []string{testRBSVolumeID, "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_credentials.json")),
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().GetCredentials(gomock.Any(), testRBSVolumeID).Return(&testRBSCredentials, nil)
			},
		},
		{
			name:        "get credentials error",
			args:        []string{testRBSVolumeID},
			expectError: true,
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().GetCredentials(gomock.Any(), testRBSVolumeID).Return(nil, errors.New("not found"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rbsVolumeServiceHandler := mocks.NewMockRemoteBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.RemoteBlockStorageVolumes = rbsVolumeServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(rbsVolumeServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rbsVolumeCmd := NewCmd(testCmdContext)

			args := []string{"rbs", "get-credentials"}
			args = append(args, tc.args...)

			builder := testutils.NewTestCommandBuilder().WithCommand(rbsVolumeCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}

func TestResetRBSVolumeCredentialsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		configureMock  func(*mocks.MockRemoteBlockStorageVolumesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "reset credentials json",
			args:           []string{testRBSVolumeID, "--output", "json"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().ResetCredentials(gomock.Any(), testRBSVolumeID).Return(&testRBSVolume, nil)
			},
		},
		{
			name:        "reset credentials error",
			args:        []string{testRBSVolumeID},
			expectError: true,
			configureMock: func(mock *mocks.MockRemoteBlockStorageVolumesService) {
				mock.EXPECT().ResetCredentials(gomock.Any(), testRBSVolumeID).Return(nil, errors.New("reset error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rbsVolumeServiceHandler := mocks.NewMockRemoteBlockStorageVolumesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.RemoteBlockStorageVolumes = rbsVolumeServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(rbsVolumeServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			rbsVolumeCmd := NewCmd(testCmdContext)

			args := []string{"rbs", "reset-credentials"}
			args = append(args, tc.args...)

			builder := testutils.NewTestCommandBuilder().WithCommand(rbsVolumeCmd).WithArgs(args)
			err := builder.Build().Execute()

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(builder.GetOutput()).To(MatchJSON(tc.expectedOutput))
			}
		})
	}
}
