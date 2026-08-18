package metrics

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBuildRackRows(t *testing.T) {
	g := NewWithT(t)

	samples, err := Parse(readFixture(t, "racks_input.txt"))
	g.Expect(err).To(BeNil())

	rows := BuildRackRows(samples)

	g.Expect(rows).To(Equal([]RackMetric{
		{
			RackID:          "0pEOrzdl",
			Title:           "rack-a",
			LocationID:      "52",
			LocationCode:    "LUX3",
			Hosts:           4,
			PublicSent:      1319413953331,
			PublicReceived:  3775348762345,
			PrivateSent:     145497332,
			PrivateReceived: 619636079,
			TotalSent:       1319559450663,
			TotalReceived:   3775968398424,
			// summed over both PDUs of the rack
			PduWatts:   1240,
			PduAmperes: 5.6,
			PduCount:   2,
			AtsWatts:   1240,
			AtsAmperes: 5.6,
			AtsCount:   1,
		},
		{
			// a rack without hosts, traffic and power devices is still listed
			RackID:       "7xKLmnQp",
			Title:        "rack-b",
			LocationID:   "52",
			LocationCode: "LUX3",
		},
	}))
}

func TestBuildRackRowsWithoutRackID(t *testing.T) {
	g := NewWithT(t)

	// the racks count metric has no rack id and produces no rows
	samples, err := Parse(`serverscom_racks_count{location_id="52",location_code="LUX3"} 2`)
	g.Expect(err).To(BeNil())

	g.Expect(BuildRackRows(samples)).To(BeEmpty())
}
