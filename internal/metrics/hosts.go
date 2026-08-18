package metrics

import (
	"cmp"
	"slices"
	"strings"
)

const (
	hostSentMetric     = "serverscom_host_monthly_sent_bytes_total"
	hostReceivedMetric = "serverscom_host_monthly_received_bytes_total"
)

// HostMetric represents metrics of a single host.
type HostMetric struct {
	HostID          string
	Title           string
	HostType        string
	ChassisName     string
	LocationID      string
	LocationCode    string
	RackID          string
	RackType        string
	PublicSent      int64
	PublicReceived  int64
	PrivateSent     int64
	PrivateReceived int64
	TotalSent       int64
	TotalReceived   int64
}

// BuildHostRows folds hosts metrics samples into one row per host.
// Only hosts with traffic data are labeled with a host id, so hosts without
// any traffic counter can't be represented as a row.
func BuildHostRows(samples []Sample) []HostMetric {
	rows := make(map[string]*HostMetric)

	for _, sample := range samples {
		if sample.Name != hostSentMetric && sample.Name != hostReceivedMetric {
			continue
		}

		id := sample.Labels["host_id"]
		if id == "" {
			continue
		}

		row, ok := rows[id]
		if !ok {
			row = &HostMetric{
				HostID:       id,
				Title:        sample.Labels["title"],
				HostType:     sample.Labels["host_type"],
				ChassisName:  sample.Labels["chassis_name"],
				LocationID:   sample.Labels["location_id"],
				LocationCode: sample.Labels["location_code"],
				RackID:       sample.Labels["rack_id"],
				RackType:     sample.Labels["rack_type"],
			}
			rows[id] = row
		}

		value := int64(sample.Value)
		sent := sample.Name == hostSentMetric

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
	}

	result := make([]HostMetric, 0, len(rows))
	for _, row := range rows {
		result = append(result, *row)
	}
	slices.SortFunc(result, func(a, b HostMetric) int {
		return cmp.Or(
			strings.Compare(a.LocationCode, b.LocationCode),
			strings.Compare(a.Title, b.Title),
			strings.Compare(a.HostID, b.HostID),
		)
	})

	return result
}
