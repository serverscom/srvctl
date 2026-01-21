package sshkeys

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
	testFingerprint      = "00:11:22:33:44:55:66:77:88:99"
	fixtureBasePath      = filepath.Join("..", "..", "..", "testdata", "entities", "ssh-keys")
	skeletonTemplatePath = filepath.Join("..", "..", "..", "internal", "output", "skeletons", "skeleton-templates", "ssh-keys")
	fixedTime            = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	testSSHKey           = serverscom.SSHKey{
		Name:        "test-key",
		Fingerprint: testFingerprint,
		Labels:      map[string]string{"foo": "bar"},
		Created:     fixedTime,
		Updated:     fixedTime,
	}
)

func TestAddSSHKeysCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockSSHKeysService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create ssh key with input",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create.json")},
			configureMock: func(mock *mocks.MockSSHKeysService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.SSHKeyCreateInput{
						Name:      "test-key",
						PublicKey: "ssh-rsa AAA",
						Labels:    map[string]string{"foo": "bar"},
					}).
					Return(&testSSHKey, nil)
			},
		},
		{
			name:           "create ssh key",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args: []string{
				"--name", "test-key",
				"--public-key", "-----TEST public-key-----",
				"--label", "foo=bar",
			},
			configureMock: func(mock *mocks.MockSSHKeysService) {
				mock.EXPECT().
					Create(gomock.Any(), serverscom.SSHKeyCreateInput{
						Name:      "test-key",
						PublicKey: "-----TEST public-key-----",
						Labels:    map[string]string{"foo": "bar"},
					}).
					Return(&testSSHKey, nil)
			},
		},
		{
			name:           "skeleton for ssh key input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add.json")),
			configureMock: func(mock *mocks.MockSSHKeysService) {
				mock.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:        "create ssh key with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sshServiceHandler := mocks.NewMockSSHKeysService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSHKeys = sshServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(sshServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"ssh-keys", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestGetSSHKeysCmd(t *testing.T) {
	testCases := []struct {
		name           string
		fingerprint    string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get ssh key in default format",
			fingerprint:    testFingerprint,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get ssh key in JSON format",
			fingerprint:    testFingerprint,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get ssh key in YAML format",
			fingerprint:    testFingerprint,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get ssh key with error",
			fingerprint: testFingerprint,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sshServiceHandler := mocks.NewMockSSHKeysService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSHKeys = sshServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			sshServiceHandler.EXPECT().
				Get(gomock.Any(), testFingerprint).
				Return(&testSSHKey, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"ssh-keys", "get", tc.fingerprint}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestListSSHKeysCmd(t *testing.T) {
	testKey1 := testSSHKey
	testKey2 := testSSHKey
	testKey2.Name = "test-key 2"
	testKey2.Fingerprint = "00:00:00:00:00:00:00:00:00:00"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.SSHKey])
	}{
		{
			name:           "list all ssh keys",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSHKey]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.SSHKey{
						testKey1,
						testKey2,
					}, nil)
			},
		},
		{
			name:           "list ssh keys",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSHKey]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SSHKey{
						testKey1,
					}, nil)
			},
		},
		{
			name:           "list ssh keys with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSHKey]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SSHKey{
						testKey1,
						testKey2,
					}, nil)
			},
		},
		{
			name:           "list ssh keys with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSHKey]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SSHKey{
						testKey1,
						testKey2,
					}, nil)
			},
		},
		{
			name:        "list ssh keys with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.SSHKey]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sshServiceHandler := mocks.NewMockSSHKeysService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.SSHKey](mockCtrl)

	sshServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSHKeys = sshServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"ssh-keys", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestUpdateSSHKeysCmd(t *testing.T) {
	newSSHKey := testSSHKey
	newSSHKey.Name = "new-ssh-key"
	newSSHKey.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		fingerprint    string
		output         string
		args           []string
		configureMock  func(*mocks.MockSSHKeysService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update ssh key",
			fingerprint:    testFingerprint,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update.json")),
			args:           []string{"--name", newSSHKey.Name, "--label", "new=label"},
			configureMock: func(mock *mocks.MockSSHKeysService) {
				mock.EXPECT().
					Update(gomock.Any(), testFingerprint, serverscom.SSHKeyUpdateInput{
						Name:   "new-ssh-key",
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newSSHKey, nil)
			},
		},
		{
			name:        "update ssh key with error",
			fingerprint: testFingerprint,
			configureMock: func(mock *mocks.MockSSHKeysService) {
				mock.EXPECT().
					Update(gomock.Any(), testFingerprint, serverscom.SSHKeyUpdateInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sshServiceHandler := mocks.NewMockSSHKeysService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSHKeys = sshServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(sshServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"ssh-keys", "update", tc.fingerprint}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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

func TestDeleteSSHKeysCmd(t *testing.T) {
	testCases := []struct {
		name        string
		fingerprint string
		expectError bool
	}{
		{
			name:        "delete ssh key",
			fingerprint: testFingerprint,
		},
		{
			name:        "delete ssh key with error",
			fingerprint: testFingerprint,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sshServiceHandler := mocks.NewMockSSHKeysService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSHKeys = sshServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			sshServiceHandler.EXPECT().
				Delete(gomock.Any(), testFingerprint).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sshCmd := NewCmd(testCmdContext)

			args := []string{"ssh-keys", "delete", tc.fingerprint}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sshCmd).
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
