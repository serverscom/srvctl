package cloudinstances

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
	ptrFixtureBasePath = filepath.Join("..", "..", "..", "testdata", "entities", "ptr")
	testInstanceID     = "instanceId"
	testPTRID          = "ptrId"

	testPTR = serverscom.PTRRecord{
		ID:       testPTRID,
		IP:       "192.0.2.5",
		Domain:   "ptr-test.example.com",
		Priority: 10,
		TTL:      300,
	}
)

func TestListPTRCmd(t *testing.T) {
	testPTR1 := testPTR
	testPTR2 := testPTR
	testPTR1.ID += "1"
	testPTR2.ID += "2"
	testPTR2.Domain = "another.example.com"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.PTRRecord])
	}{
		{
			name:           "list all ptr records",
			output:         "json",
			args:           []string{"-A", testInstanceID},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:           "list ptr records",
			output:         "json",
			args:           []string{testInstanceID},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
					}, nil)
			},
		},
		{
			name:           "list ptr records with template",
			args:           []string{testInstanceID, "--template", "{{range .}}ID: {{.ID}} PTR: {{.Domain}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:           "list ptr records with page-view",
			args:           []string{testInstanceID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.PTRRecord{
						testPTR1,
						testPTR2,
					}, nil)
			},
		},
		{
			name:        "list ptr records with error",
			args:        []string{testInstanceID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.PTRRecord]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloudInstanceService := mocks.NewMockCloudComputingInstancesService(mockCtrl)
	collection := mocks.NewMockCollection[serverscom.PTRRecord](mockCtrl)

	cloudInstanceService.EXPECT().PTRRecords(gomock.Any()).Return(collection).AnyTimes()
	collection.EXPECT().SetParam(gomock.Any(), gomock.Any()).Return(collection).AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudInstanceService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(collection)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := append([]string{"cloud-instances", "list-ptr"}, tc.args...)
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

func TestAddPTRCmd(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		output         string
		expectedOutput []byte
		configureMock  func(*mocks.MockCloudComputingInstancesService)
		expectError    bool
	}{
		{
			name:           "add ptr record",
			args:           []string{testInstanceID, "--data", testPTR.Domain, "--ip", testPTR.IP, "--ttl", "300", "--priority", "10"},
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(ptrFixtureBasePath, "get.json")),
			configureMock: func(mock *mocks.MockCloudComputingInstancesService) {
				in := serverscom.PTRRecordCreateInput{
					Domain:   testPTR.Domain,
					IP:       testPTR.IP,
					TTL:      &testPTR.TTL,
					Priority: &testPTR.Priority,
				}
				mock.EXPECT().
					CreatePTRRecord(gomock.Any(), testInstanceID, in).
					Return(&testPTR, nil)
			},
		},
		{
			name:        "add ptr record error",
			args:        []string{testInstanceID, "--data", testPTR.Domain, "--ip", testPTR.IP},
			expectError: true,
			configureMock: func(mock *mocks.MockCloudComputingInstancesService) {
				in := serverscom.PTRRecordCreateInput{
					Domain: testPTR.Domain,
					IP:     testPTR.IP,
				}
				mock.EXPECT().
					CreatePTRRecord(gomock.Any(), testInstanceID, in).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloudInstanceService := mocks.NewMockCloudComputingInstancesService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudInstanceService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(cloudInstanceService)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "add-ptr"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
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

func TestDeletePTRCmd(t *testing.T) {
	testCases := []struct {
		name        string
		instanceID  string
		ptrID       string
		expectError bool
	}{
		{
			name:       "delete ptr record",
			instanceID: testInstanceID,
			ptrID:      testPTRID,
		},
		{
			name:        "delete ptr record error",
			instanceID:  testInstanceID,
			ptrID:       testPTRID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloudInstanceService := mocks.NewMockCloudComputingInstancesService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.CloudComputingInstances = cloudInstanceService

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			var err error
			if tc.expectError {
				err = errors.New("some error")
			}
			cloudInstanceService.EXPECT().
				DeletePTRRecord(gomock.Any(), tc.instanceID, tc.ptrID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			cmd := NewCmd(testCmdContext)

			args := []string{"cloud-instances", "delete-ptr", tc.instanceID, "--ptr-id", tc.ptrID}
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
