package metrics

import (
	"cmp"
	"slices"
	"strings"
)

const (
	rackHostsCountMetric = "serverscom_rack_hosts_count"
	rackSentMetric       = "serverscom_rack_monthly_sent_bytes_total"
	rackReceivedMetric   = "serverscom_rack_monthly_received_bytes_total"
	rackPduPowerMetric   = "serverscom_rack_pdu_power_watts"
	rackPduCurrentMetric = "serverscom_rack_pdu_current_amperes"
	rackAtsPowerMetric   = "serverscom_rack_ats_power_watts"
	rackAtsCurrentMetric = "serverscom_rack_ats_current_amperes"
)

// RackMetric represents metrics of a single rack.
// Power and current are summed per device type. PDU and ATS are kept apart
// because an ATS feeds the PDUs, so summing them would count the same draw twice.
type RackMetric struct {
	RackID          string
	Title           string
	LocationID      string
	LocationCode    string
	Hosts           int64
	PublicSent      int64
	PublicReceived  int64
	PrivateSent     int64
	PrivateReceived int64
	TotalSent       int64
	TotalReceived   int64
	PduWatts        float64
	PduAmperes      float64
	PduCount        int
	AtsWatts        float64
	AtsAmperes      float64
	AtsCount        int
}

// BuildRackRows folds racks metrics samples into one row per rack.
func BuildRackRows(samples []Sample) []RackMetric {
	rows := make(map[string]*RackMetric)
	// power and current are reported per device, count each device once
	devices := make(map[string]map[string]struct{})

	getRow := func(sample Sample) *RackMetric {
		id := sample.Labels["rack_id"]
		if id == "" {
			return nil
		}
		row, ok := rows[id]
		if !ok {
			row = &RackMetric{
				RackID:       id,
				Title:        sample.Labels["rack_title"],
				LocationID:   sample.Labels["location_id"],
				LocationCode: sample.Labels["location_code"],
			}
			rows[id] = row
		}
		return row
	}

	countDevice := func(id, deviceType, name string) bool {
		key := deviceType + "/" + id
		if devices[key] == nil {
			devices[key] = make(map[string]struct{})
		}
		if _, ok := devices[key][name]; ok {
			return false
		}
		devices[key][name] = struct{}{}
		return true
	}

	for _, sample := range samples {
		row := getRow(sample)
		if row == nil {
			continue
		}

		switch sample.Name {
		case rackHostsCountMetric:
			row.Hosts += int64(sample.Value)
		case rackSentMetric, rackReceivedMetric:
			value := int64(sample.Value)
			sent := sample.Name == rackSentMetric

			switch sample.Labels["traffic_type"] {
			case "public":
				if sent {
					row.PublicSent += value
				} else {
					row.PublicReceived += value
				}
			case "private":
				if sent {
					row.PrivateSent += value
				} else {
					row.PrivateReceived += value
				}
			}

			// totals cover all traffic types, including the ones without a column
			if sent {
				row.TotalSent += value
			} else {
				row.TotalReceived += value
			}
		case rackPduPowerMetric:
			row.PduWatts += sample.Value
			if countDevice(row.RackID, "pdu", sample.Labels["pdu_name"]) {
				row.PduCount++
			}
		case rackPduCurrentMetric:
			row.PduAmperes += sample.Value
			if countDevice(row.RackID, "pdu", sample.Labels["pdu_name"]) {
				row.PduCount++
			}
		case rackAtsPowerMetric:
			row.AtsWatts += sample.Value
			if countDevice(row.RackID, "ats", sample.Labels["ats_name"]) {
				row.AtsCount++
			}
		case rackAtsCurrentMetric:
			row.AtsAmperes += sample.Value
			if countDevice(row.RackID, "ats", sample.Labels["ats_name"]) {
				row.AtsCount++
			}
		}
	}

	result := make([]RackMetric, 0, len(rows))
	for _, row := range rows {
		result = append(result, *row)
	}
	slices.SortFunc(result, func(a, b RackMetric) int {
		return cmp.Or(
			strings.Compare(a.LocationCode, b.LocationCode),
			strings.Compare(a.Title, b.Title),
			strings.Compare(a.RackID, b.RackID),
		)
	})

	return result
}
