package l2segments

import (
	"errors"
	"fmt"
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
	fixtureBasePath      = filepath.Join("..", "..", "..", "testdata", "entities", "l2-segments")
	skeletonTemplatePath = filepath.Join("..", "..", "..", "internal", "output", "skeletons", "templates", "l2-segments")
	testID               = "testId"
	testL2SegmentName    = "testName"
	fixedTime            = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	testL2Segment = serverscom.L2Segment{
		ID:                testID,
		Name:              testL2SegmentName,
		Type:              "public",
		Status:            "pending",
		LocationGroupID:   58,
		LocationGroupCode: "location1565",
		Labels: map[string]string{
			"environment": "testing",
		},
		Created: fixedTime,
		Updated: fixedTime,
	}

	testL2LocationGroup = serverscom.L2LocationGroup{
		ID:           10,
		Name:         testL2SegmentName,
		Code:         "testCode",
		GroupType:    "public",
		LocationIDs:  []int64{1, 2, 3},
		Hyperscalers: []string{"AWS", "Azure"},
	}

	testL2Member = serverscom.L2Member{
		ID:     testID,
		Title:  "member-title",
		Mode:   "native",
		Status: "new",
		Labels: map[string]string{
			"environment": "testing",
		},
		Created: fixedTime,
		Updated: fixedTime,
	}

	cidr        = "127.1.182.0/24"
	title       = "testNetwork"
	testNetwork = serverscom.Network{
		ID:                 testID,
		Title:              &title,
		Status:             "new",
		Cidr:               &cidr,
		Family:             "ipv4",
		InterfaceType:      "public",
		DistributionMethod: "route",
		Additional:         true,
		Created:            fixedTime,
		Updated:            fixedTime,
	}
)

func TestGetL2SegmentCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "get l2 segment in default format",
			id:             testID,
			output:         "",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.txt")),
		},
		{
			name:           "get l2 segment in JSON format",
			id:             testID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
		},
		{
			name:           "get l2 segment in YAML format",
			id:             testID,
			output:         "yaml",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.yaml")),
		},
		{
			name:        "get l2 segment with error",
			id:          testID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}

			l2ServiceHandler.EXPECT().
				Get(gomock.Any(), testID).
				Return(&testL2Segment, err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "get", fmt.Sprint(tc.id)}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestListL2SegmentsCmd(t *testing.T) {
	testL2Segment1 := testL2Segment
	testL2Segment2 := testL2Segment
	testL2Segment1.ID += "1"
	testL2Segment2.Name = "other-segment"
	testL2Segment2.ID += "2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.L2Segment])
	}{
		{
			name:           "list all l2 segments",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Segment]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.L2Segment{
						testL2Segment1,
						testL2Segment2,
					}, nil)
			},
		},
		{
			name:           "list l2 segments",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Segment]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2Segment{
						testL2Segment1,
					}, nil)
			},
		},
		{
			name:           "list l2 segments with template",
			args:           []string{"--template", "{{range .}}Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Segment]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2Segment{
						testL2Segment1,
						testL2Segment2,
					}, nil)
			},
		},
		{
			name:           "list l2 segments with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Segment]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2Segment{
						testL2Segment1,
						testL2Segment2,
					}, nil)
			},
		},
		{
			name:        "list l2 segments with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Segment]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	collectionHandler := mocks.NewMockCollection[serverscom.L2Segment](mockCtrl)

	l2ServiceHandler.EXPECT().
		Collection().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "list"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestAddL2SegmentCmd(t *testing.T) {
	testCases := []struct {
		name           string
		output         string
		args           []string
		configureMock  func(*mocks.MockL2SegmentsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "create l2 segment with input",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "create.json")},
			configureMock: func(mock *mocks.MockL2SegmentsService) {
				mock.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(serverscom.L2SegmentCreateInput{})).
					Return(&testL2Segment, nil)
			},
		},
		{
			name:           "create l2 segment",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args: []string{
				"--name", testL2SegmentName,
				"--type", "public",
				"--location-group-id", "58",
				"--member", "id=LDdwmwa1,mode=native",
				"--member", "id=LDdwmwa2,mode=trunk",
				"--member", "id=LDdwmwa3,mode=native",
				"--label", "foo=bar",
				"--label", "bar=foo",
			},
			configureMock: func(mock *mocks.MockL2SegmentsService) {
				expectedMembers := []serverscom.L2SegmentMemberInput{
					{ID: "LDdwmwa1", Mode: "native"},
					{ID: "LDdwmwa2", Mode: "trunk"},
					{ID: "LDdwmwa3", Mode: "native"},
				}
				expectedLabels := map[string]string{
					"foo": "bar",
					"bar": "foo",
				}

				mock.EXPECT().
					Create(gomock.Any(), serverscom.L2SegmentCreateInput{
						Name:            &testL2SegmentName,
						Type:            "public",
						LocationGroupID: 58,
						Members:         expectedMembers,
						Labels:          expectedLabels,
					}).
					Return(&testL2Segment, nil)
			},
		},
		{
			name:           "skeleton for l2 segment input",
			output:         "json",
			args:           []string{"--skeleton"},
			expectedOutput: testutils.ReadFixture(filepath.Join(skeletonTemplatePath, "add.json")),
			configureMock: func(mock *mocks.MockL2SegmentsService) {
				mock.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name:        "create l2 segment with error",
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(l2ServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "add"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestUpdateL2SegmentCmd(t *testing.T) {
	updatedSeg := testL2Segment
	updatedSeg.Labels = map[string]string{"new": "label"}

	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockL2SegmentsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update l2 segment",
			id:             testID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "update.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "update_input.json")},
			configureMock: func(mock *mocks.MockL2SegmentsService) {
				mock.EXPECT().
					Update(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.L2SegmentUpdateInput{})).
					Return(&updatedSeg, nil)
			},
		},
		{
			name: "update l2 segment with error",
			id:   testID,
			args: []string{"--input", filepath.Join(fixtureBasePath, "update_input.json")},
			configureMock: func(mock *mocks.MockL2SegmentsService) {
				mock.EXPECT().
					Update(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.L2SegmentUpdateInput{})).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(l2ServiceHandler)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "update", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestDeleteL2SegmentCmd(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name: "delete l2 segment",
			id:   testID,
		},
		{
			name:        "delete l2 segment with error",
			id:          testID,
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			var err error
			if tc.expectError {
				err = errors.New("some error")
			}

			l2ServiceHandler.EXPECT().
				Delete(gomock.Any(), testID).
				Return(err)

			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "delete", tc.id}
			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestUpdateL2NetworksCmd(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		output         string
		args           []string
		configureMock  func(*mocks.MockL2SegmentsService)
		expectedOutput []byte
		expectError    bool
	}{
		{
			name:           "update l2 segment networks",
			id:             testID,
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "get.json")),
			args:           []string{"--input", filepath.Join(fixtureBasePath, "update_networks_input.json")},
			configureMock: func(mock *mocks.MockL2SegmentsService) {
				mock.EXPECT().
					ChangeNetworks(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.L2SegmentChangeNetworksInput{})).
					Return(&testL2Segment, nil)
			},
		},
		{
			name: "update l2 segment networks with error",
			id:   testID,
			args: []string{"--input", filepath.Join(fixtureBasePath, "update_networks_input.json")},
			configureMock: func(mock *mocks.MockL2SegmentsService) {
				mock.EXPECT().
					ChangeNetworks(gomock.Any(), testID, gomock.AssignableToTypeOf(serverscom.L2SegmentChangeNetworksInput{})).
					Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(l2ServiceHandler)
			}

			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "update-networks", tc.id}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestListL2GroupsCmd(t *testing.T) {
	g1 := testL2LocationGroup
	g2 := testL2LocationGroup
	g2.ID += 1
	g2.Name += "2"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.L2LocationGroup])
	}{
		{
			name:           "list all l2 groups",
			output:         "json",
			args:           []string{"-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_groups_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2LocationGroup]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.L2LocationGroup{g1, g2}, nil)
			},
		},
		{
			name:           "list l2 groups",
			output:         "json",
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_groups.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2LocationGroup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2LocationGroup{g1}, nil)
			},
		},
		{
			name:           "list l2 groups with template",
			args:           []string{"--template", "{{range .}}ID: {{.ID}} Name: {{.Name}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_groups_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2LocationGroup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2LocationGroup{g1, g2}, nil)
			},
		},
		{
			name:           "list l2 groups with pageView",
			args:           []string{"--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_groups_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2LocationGroup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2LocationGroup{g1, g2}, nil)
			},
		},
		{
			name:        "list l2 groups with error",
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.L2LocationGroup]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	collectionHandler := mocks.NewMockCollection[serverscom.L2LocationGroup](mockCtrl)
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	l2ServiceHandler.EXPECT().
		LocationGroups().
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "list-groups"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestListL2MembersCmd(t *testing.T) {
	m1 := testL2Member
	m2 := testL2Member
	m1.ID += "1"
	m2.ID += "2"
	vlan := 100
	m2.Vlan = &vlan
	m2.Title = "other-member"

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.L2Member])
	}{
		{
			name:           "list all l2 members",
			output:         "json",
			args:           []string{testID, "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_members_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Member]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.L2Member{m1, m2}, nil)
			},
		},
		{
			name:           "list l2 members",
			output:         "json",
			args:           []string{testID},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_members.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Member]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2Member{m1}, nil)
			},
		},
		{
			name:           "list l2 members with template",
			args:           []string{testID, "--template", "{{range .}}ID: {{.ID}} Title: {{.Title}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_members_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Member]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2Member{m1, m2}, nil)
			},
		},
		{
			name:           "list l2 members with pageView",
			args:           []string{testID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_members_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Member]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.L2Member{m1, m2}, nil)
			},
		},
		{
			name:        "list l2 members error",
			args:        []string{testID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.L2Member]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	collectionHandler := mocks.NewMockCollection[serverscom.L2Member](mockCtrl)
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	l2ServiceHandler.EXPECT().
		Members(testID).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "list-members"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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

func TestListL2NetworksCmd(t *testing.T) {
	n1 := testNetwork
	n2 := testNetwork
	n1.ID += "1"
	n2.ID += "2"
	title := "other-network"
	n2.Title = &title

	testCases := []struct {
		name           string
		output         string
		args           []string
		expectedOutput []byte
		expectError    bool
		configureMock  func(*mocks.MockCollection[serverscom.Network])
	}{
		{
			name:           "list all l2 networks",
			output:         "json",
			args:           []string{testID, "-A"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_networks_all.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					Collect(gomock.Any()).
					Return([]serverscom.Network{n1, n2}, nil)
			},
		},
		{
			name:           "list l2 networks",
			output:         "json",
			args:           []string{testID},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_networks.json")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{n1}, nil)
			},
		},
		{
			name:           "list l2 networks with template",
			args:           []string{testID, "--template", "{{range .}}ID: {{.ID}} Title: {{.Title}}\n{{end}}"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_networks_template.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{n1, n2}, nil)
			},
		},
		{
			name:           "list l2 networks with pageView",
			args:           []string{testID, "--page-view"},
			expectedOutput: testutils.ReadFixture(filepath.Join(fixtureBasePath, "list_networks_pageview.txt")),
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return([]serverscom.Network{n1, n2}, nil)
			},
		},
		{
			name:        "list l2 networks error",
			args:        []string{testID},
			expectError: true,
			configureMock: func(mock *mocks.MockCollection[serverscom.Network]) {
				mock.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("some error"))
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	collectionHandler := mocks.NewMockCollection[serverscom.Network](mockCtrl)
	l2ServiceHandler := mocks.NewMockL2SegmentsService(mockCtrl)
	l2ServiceHandler.EXPECT().
		Networks(testID).
		Return(collectionHandler).
		AnyTimes()

	collectionHandler.EXPECT().
		SetParam(gomock.Any(), gomock.Any()).
		Return(collectionHandler).
		AnyTimes()

	scClient := serverscom.NewClientWithEndpoint("", "")
	scClient.L2Segments = l2ServiceHandler

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			if tc.configureMock != nil {
				tc.configureMock(collectionHandler)
			}
			testCmdContext := testutils.NewTestCmdContext(scClient)
			l2Cmd := NewCmd(testCmdContext)

			args := []string{"l2-segments", "list-networks"}
			if len(tc.args) > 0 {
				args = append(args, tc.args...)
			}
			if tc.output != "" {
				args = append(args, "--output", tc.output)
			}

			builder := testutils.NewTestCommandBuilder().
				WithCommand(l2Cmd).
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
