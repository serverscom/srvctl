package cloudbackups

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
	testBackupID            = "backup-123"
	testVolumeUUID          = "vol-123-openstack"
	testVolumeUUID2         = "vol-456-openstack"
	fixtureBasePath         = filepath.Join("..", "..", "..", "testdata", "entities", "cloud-backups")
	fixedTime               = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	fixedTime2              = time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC)
	testBackupOpenstackUUID *string
	testBackup              = serverscom.CloudBlockStorageBackup{
		ID:                  testBackupID,
		OpenstackUUID:       testBackupOpenstackUUID,
		OpenstackVolumeUUID: testVolumeUUID,
		RegionID:            1,
		Size:                1073741824,
		Status:              "available",
		Labels:              map[string]string{"env": "test"},
		Created:             &fixedTime,
		Name:                "test-backup",
	}
	testBackup2 = serverscom.CloudBlockStorageBackup{
		ID:                  "backup-456",
		OpenstackUUID:       testBackupOpenstackUUID,
		OpenstackVolumeUUID: testVolumeUUID2,
		RegionID:            1,
		Size:                2147483648,
		Status:              "available",
		Labels:              map[string]string{"env": "prod"},
		Created:             &fixedTime2,
		Name:                "test-backup-2",
	}
)

func TestAddCloudBackupsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageBackupsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create cloud backup with flags",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args: []string{
				"--volume-id", testVolumeUUID,
				"--name", "test-backup",
				"--label", "env=test",
			},
			configureMock: func(mock *mocks.MockCloudBlockStorageBackupsService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.CloudBlockStorageBackupCreateInput{
						VolumeID:    testVolumeUUID,
						Name:        "test-backup",
						Incremental: false,
						Force:       false,
						Labels:      map[string]string{"env": "test"},
					}).
					Return(&testBackup, nil)
			},
		},
		{
			name:        "create cloud backup with error",
			expectError: true,
			args:        []string{"--volume-id", testVolumeUUID, "--name", "test-backup"},
			configureMock: func(mock *mocks.MockCloudBlockStorageBackupsService) {
				mock.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	backupServiceHandler := mocks.NewMockCloudBlockStorageBackupsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageBackups = backupServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(backupServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			backupCmd := NewCmd(testCmdContext)

			args := []string{"cloud-backups", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(backupCmd).
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

func TestGetCloudBackupsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		backupID       string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get cloud backup in default format",
			backupID:       testBackupID,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get cloud backup in JSON format",
			backupID:       testBackupID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get cloud backup in YAML format",
			backupID:       testBackupID,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get cloud backup with error",
			backupID:    testBackupID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	backupServiceHandler := mocks.NewMockCloudBlockStorageBackupsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageBackups = backupServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			backupServiceHandler.EXPECT().
				Get(gomock.Any(), testBackupID).
				Return(&testBackup, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			backupCmd := NewCmd(testCmdContext)

			args := []string{"cloud-backups", "get", tc.backupID}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(backupCmd).
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

func TestListCloudBackupsCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.CloudBlockStorageBackup])
	}{
		{
			name:           "list all cloud backups",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudBlockStorageBackup]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.CloudBlockStorageBackup{
						testBackup,
						testBackup2,
					}, nil)
			},
		},
		{
			name:           "list cloud backups",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudBlockStorageBackup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudBlockStorageBackup{
						testBackup,
					}, nil)
			},
		},
		{
			name:           "list cloud backups with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}}, Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudBlockStorageBackup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudBlockStorageBackup{
						testBackup,
						testBackup2,
					}, nil)
			},
		},
		{
			name:           "list cloud backups with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudBlockStorageBackup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.CloudBlockStorageBackup{
						testBackup,
						testBackup2,
					}, nil)
			},
		},
		{
			name:        "list cloud backups with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.CloudBlockStorageBackup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	backupServiceHandler := mocks.NewMockCloudBlockStorageBackupsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.CloudBlockStorageBackup](mockCtrl)

	backupServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageBackups = backupServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			backupCmd := NewCmd(testCmdContext)

			args := []string{"cloud-backups", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(backupCmd).
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

func TestUpdateCloudBackupsCmd(t *testing.T) {
	updatedBackup := testBackup
	updatedBackup.Labels = map[string]string{"env": "updated"}

	testCases := []struct {
		name           string
		backupID       string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageBackupsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update cloud backup",
			backupID:       testBackupID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update.json")),
			args:           []string{"--label", "env=updated"},
			configureMock: func(mock *mocks.MockCloudBlockStorageBackupsService) {
				mock.EXPECT().
					Update(gomock.Any(), testBackupID, serverscom.CloudBlockStorageBackupUpdateInput{
						Labels: map[string]string{"env": "updated"},
					}).
					Return(&updatedBackup, nil)
			},
		},
		{
			name:        "update cloud backup with error",
			backupID:    testBackupID,
			expectError: true,
			configureMock: func(mock *mocks.MockCloudBlockStorageBackupsService) {
				mock.EXPECT().
					Update(gomock.Any(), testBackupID, serverscom.CloudBlockStorageBackupUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	backupServiceHandler := mocks.NewMockCloudBlockStorageBackupsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageBackups = backupServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(backupServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			backupCmd := NewCmd(testCmdContext)

			args := []string{"cloud-backups", "update", tc.backupID}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(backupCmd).
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

func TestDeleteCloudBackupsCmd(t *testing.T) {
	testCases := []struct {
		name        string
		backupID    string
		expectError bool
	}{
		{
			name:     "delete cloud backup",
			backupID: testBackupID,
		},
		{
			name:        "delete cloud backup with error",
			backupID:    testBackupID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	backupServiceHandler := mocks.NewMockCloudBlockStorageBackupsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageBackups = backupServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			backupServiceHandler.EXPECT().
				Delete(gomock.Any(), testBackupID).
				Return(&testBackup, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			backupCmd := NewCmd(testCmdContext)

			args := []string{"cloud-backups", "delete", tc.backupID}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(backupCmd).
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

func TestRestoreCloudBackupsCmd(t *testing.T) {
	restoredBackup := testBackup
	restoredBackup.OpenstackVolumeUUID = testVolumeUUID2

	testCases := []struct {
		name           string
		backupID       string
		output         string
		args           []string
		configureMock  func(*mocks.MockCloudBlockStorageBackupsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "restore cloud backup",
			backupID:       testBackupID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "restore.json")),
			args:           []string{"--volume-id", testVolumeUUID2},
			configureMock: func(mock *mocks.MockCloudBlockStorageBackupsService) {
				mock.EXPECT().
					Restore(gomock.Any(), testBackupID, serverscom.CloudBlockStorageBackupRestoreInput{
						VolumeID: testVolumeUUID2,
					}).
					Return(&restoredBackup, nil)
			},
		},
		{
			name:        "restore cloud backup with error",
			backupID:    testBackupID,
			expectError: true,
			args:        []string{"--volume-id", testVolumeUUID2},
			configureMock: func(mock *mocks.MockCloudBlockStorageBackupsService) {
				mock.EXPECT().
					Restore(gomock.Any(), testBackupID, gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	backupServiceHandler := mocks.NewMockCloudBlockStorageBackupsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudBlockStorageBackups = backupServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(backupServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			backupCmd := NewCmd(testCmdContext)

			args := []string{"cloud-backups", "restore", tc.backupID}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(backupCmd).
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
