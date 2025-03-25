package ssl

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
	"github.com/serverscom/srvctl/cmd/testutils"
	"github.com/serverscom/srvctl/internal/mocks"

	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var (
	fixtureBasePath     = filepath.Join("..", "..", "..", "testdata", "entities", "ssl")
	fixedTime           = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	expiresTime         = fixedTime.AddDate(10, 0, 0)
	testSha1Fingerpring = "21e84c9a3878673b377f0adf053290e8fc25cb80"
	testIssuer          = "servers.com"
	testId              = "testId"
	testSSL             = serverscom.SSLCertificate{
		ID:              testId,
		Name:            "test-ssl-custom",
		Type:            "custom",
		Issuer:          &testIssuer,
		Subject:         "servers.com",
		DomainNames:     []string{"servers.com"},
		Sha1Fingerprint: testSha1Fingerpring,
		Labels:          map[string]string{"foo": "bar"},
		Expires:         &expiresTime,
		Created:         fixedTime,
		Updated:         fixedTime,
	}
	testCustomSSL = serverscom.SSLCertificateCustom{
		ID:              testId,
		Name:            "test-ssl-custom",
		Type:            "custom",
		Issuer:          &testIssuer,
		Subject:         "servers.com",
		DomainNames:     []string{"servers.com"},
		Sha1Fingerprint: testSha1Fingerpring,
		Labels:          map[string]string{"foo": "bar"},
		Expires:         &expiresTime,
		Created:         fixedTime,
		Updated:         fixedTime,
	}
	testLeSSL = serverscom.SSLCertificateLE{
		ID:          testId,
		Name:        "test-ssl-le",
		Type:        "letsencrypt",
		Issuer:      &testIssuer,
		Subject:     "servers.com",
		DomainNames: []string{"servers.com"},
		Labels:      map[string]string{"foo": "bar"},
		Expires:     &expiresTime,
		Created:     fixedTime,
		Updated:     fixedTime,
	}
)

func TestAddCustomSSLCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockSSLCertificatesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create custom ssl cert",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_custom.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create_custom.json")},
			configureMock: func(mock *mocks.MockSSLCertificatesService) {
				mock.EXPECT().
					CreateCustom(gomock.Any(), serverscom.SSLCertificateCreateCustomInput{
						Name:       "test-ssl-custom",
						PublicKey:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
						Labels:     map[string]string{"foo": "bar"},
					}).
					Return(&testCustomSSL, nil)
			},
		},
		{
			name:        "with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(sslServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "custom", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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

func TestGetCustomSSLCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get custom ssl cert in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_custom.txt")),
		},
		{
			name:           "get custom ssl cert JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_custom.json")),
		},
		{
			name:           "get custom ssl cert YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_custom.yaml")),
		},
		{
			name:        "get custom ssl cert with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			sslServiceHandler.EXPECT().
				GetCustom(gomock.Any(), testId).
				Return(&testCustomSSL, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "custom", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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

func TestGetLeSSLCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get le ssl cert in default format",
			id:             testId,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_le.txt")),
		},
		{
			name:           "get le ssl cert JSON format",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_le.json")),
		},
		{
			name:           "get le ssl cert YAML format",
			id:             testId,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get_le.yaml")),
		},
		{
			name:        "get le ssl cert with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			sslServiceHandler.EXPECT().
				GetLE(gomock.Any(), testId).
				Return(&testLeSSL, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "le", "get", tc.id}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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

func TestListSSLCmd(t *testing.T) {
	ssl1 := testSSL
	ssl2 := testSSL
	ssl2.Name = "test-ssl-le"
	ssl2.Type = "letsencrypt"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.SSLCertificate])
	}{
		{
			name:           "list all ssl certs",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSLCertificate]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.SSLCertificate{
						ssl1,
						ssl2,
					}, nil)
			},
		},
		{
			name:           "list custom ssl certs",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_custom.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSLCertificate]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SSLCertificate{
						ssl1,
					}, nil)
			},
		},
		{
			name:           "list le ssl certs",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_le.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSLCertificate]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SSLCertificate{
						ssl2,
					}, nil)
			},
		},
		{
			name:           "list ssl certs with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSLCertificate]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SSLCertificate{
						ssl1,
						ssl2,
					}, nil)
			},
		},
		{
			name:           "list ssl certs with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.SSLCertificate]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.SSLCertificate{
						ssl1,
						ssl2,
					}, nil)
			},
		},
		{
			name:        "list ssl certs with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.SSLCertificate]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.SSLCertificate](mockCtrl)

	sslServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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

func TestUpdateCustomSSLCmd(t *testing.T) {
	newSSL := testCustomSSL
	newSSL.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockSSLCertificatesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update ssl cert",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_custom.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockSSLCertificatesService) {
				mock.EXPECT().
					UpdateCustom(gomock.Any(), testId, serverscom.SSLCertificateUpdateCustomInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newSSL, nil)
			},
		},
		{
			name: "update ssl cert with error",
			id:   testId,
			configureMock: func(mock *mocks.MockSSLCertificatesService) {
				mock.EXPECT().
					UpdateCustom(gomock.Any(), testId, serverscom.SSLCertificateUpdateCustomInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(sslServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "custom", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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

func TestUpdateLeSSLCmd(t *testing.T) {
	newSSL := testLeSSL
	newSSL.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockSSLCertificatesService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update ssl cert",
			id:             testId,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update_le.json")),
			args:           []string{"--label", "new=label"},
			configureMock: func(mock *mocks.MockSSLCertificatesService) {
				mock.EXPECT().
					UpdateLE(gomock.Any(), testId, serverscom.SSLCertificateUpdateLEInput{
						Labels: map[string]string{"new": "label"},
					}).
					Return(&newSSL, nil)
			},
		},
		{
			name: "update ssl cert with error",
			id:   testId,
			configureMock: func(mock *mocks.MockSSLCertificatesService) {
				mock.EXPECT().
					UpdateLE(gomock.Any(), testId, serverscom.SSLCertificateUpdateLEInput{
						Labels: make(map[string]string),
					}).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(sslServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "le", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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

func TestDeleteCustomSSLKeysCmd(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name: "delete ssl cert",
			id:   testId,
		},
		{
			name:        "delete ssl cert with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			sslServiceHandler.EXPECT().
				DeleteCustom(gomock.Any(), testId).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "custom", "delete", tc.id}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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

func TestDeleteLeSSLKeysCmd(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name: "delete ssl cert",
			id:   testId,
		},
		{
			name:        "delete ssl cert with error",
			id:          testId,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sslServiceHandler := mocks.NewMockSSLCertificatesService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.SSLCertificates = sslServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			sslServiceHandler.EXPECT().
				DeleteLE(gomock.Any(), testId).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			sslCmd := NewCmd(testCmdContext)

			args := []string{"ssl", "le", "delete", tc.id}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(sslCmd).
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
