package metrics

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

var fixtureBasePath = filepath.Join("..", "..", "testdata", "entities", "metrics")

func readFixture(t *testing.T, name string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(fixtureBasePath, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestBuildHostRows(t *testing.T) {
	g := NewWithT(t)

	samples, err := Parse(readFixture(t, "hosts_input.txt"))
	g.Expect(err).To(BeNil())

	// the count metric reports 3 hosts, only the 2 with traffic data can be rows
	rows := BuildHostRows(samples)

	g.Expect(rows).To(Equal([]HostMetric{
		{
			HostID:         "5VmrzVmx",
			Title:          "lon1-web-01",
			HostType:       "dedicated_server",
			ChassisName:    `Dell R330 - E3-1230 v6 - 3.5"`,
			LocationID:     "23",
			LocationCode:   "LON1",
			RackID:         "5VmrzVmx",
			RackType:       "shared",
			PublicSent:     1319413953331,
			PublicReceived: 3775348762345,
			TotalSent:      1319413953331,
			TotalReceived:  3775348762345,
		},
		{
			HostID:          "jpAAGYJp",
			Title:           "lux3test3-reordered",
			HostType:        "dedicated_server",
			ChassisName:     `Dell R440 - Silver 4114 - 2.5"`,
			LocationID:      "52",
			LocationCode:    "LUX3",
			RackID:          "0pEOrzdl",
			RackType:        "shared",
			PublicSent:      146447194,
			PublicReceived:  540736516,
			PrivateSent:     145497332,
			PrivateReceived: 619636079,
			TotalSent:       291944526,
			TotalReceived:   1160372595,
		},
	}))
}

func TestBuildHostRowsUnknownTrafficType(t *testing.T) {
	g := NewWithT(t)

	// a traffic type without its own column still has to be counted in the totals
	samples, err := Parse(
		`serverscom_host_monthly_sent_bytes_total{host_id="a",traffic_type="public"} 100` + "\n" +
			`serverscom_host_monthly_sent_bytes_total{host_id="a",traffic_type="unknown"} 20` + "\n" +
			`serverscom_host_monthly_received_bytes_total{host_id="a",traffic_type="unknown"} 5` + "\n",
	)
	g.Expect(err).To(BeNil())

	rows := BuildHostRows(samples)

	g.Expect(rows).To(HaveLen(1))
	g.Expect(rows[0].PublicSent).To(BeEquivalentTo(100))
	g.Expect(rows[0].TotalSent).To(BeEquivalentTo(120))
	g.Expect(rows[0].TotalReceived).To(BeEquivalentTo(5))
}

func TestBuildHostRowsWithoutHostID(t *testing.T) {
	g := NewWithT(t)

	samples, err := Parse(`serverscom_host_monthly_sent_bytes_total{traffic_type="public"} 100`)
	g.Expect(err).To(BeNil())

	rows := BuildHostRows(samples)

	g.Expect(rows).To(BeEmpty())
}
