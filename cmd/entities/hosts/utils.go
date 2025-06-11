package hosts

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	serverscom "github.com/serverscom/serverscom-go-client/pkg"
)

type parsedPartition struct {
	Slots     []int
	Partition serverscom.DedicatedServerLayoutPartitionInput
}

// parseDriveSlots parses driveSlots map into []DedicatedServerSlotInput
func parseDriveSlots(driveSlots map[string]int) ([]serverscom.DedicatedServerSlotInput, error) {
	slots := make([]serverscom.DedicatedServerSlotInput, 0, len(driveSlots))
	for pos, id := range driveSlots {
		dId := int64(id)
		posNum, err := strconv.Atoi(pos)
		if err != nil {
			return nil, fmt.Errorf("can't parse drive slot position '%s' as integer", pos)
		}
		slots = append(slots, serverscom.DedicatedServerSlotInput{
			Position:     posNum,
			DriveModelID: &dId,
		})
	}
	return slots, nil
}

// parseLayout parses layouts strings into []DedicatedServerLayoutInput
func parseLayout(layouts []string) ([]serverscom.DedicatedServerLayoutInput, error) {
	var result []serverscom.DedicatedServerLayoutInput

	for _, l := range layouts {
		var lInput serverscom.DedicatedServerLayoutInput
		parts := strings.Split(l, ",")

		for _, part := range parts {
			pair := strings.SplitN(part, "=", 2)
			if len(pair) != 2 {
				continue
			}
			key := pair[0]
			val := pair[1]

			switch key {
			case "slot":
				num, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("can't parse layout slot '%s' as integer", val)
				}
				lInput.SlotPositions = append(lInput.SlotPositions, num)
			case "raid":
				raid, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("can't parse layout raid '%s' as integer", val)
				}
				lInput.Raid = &raid
			}
		}
		if len(lInput.SlotPositions) == 0 {
			return nil, fmt.Errorf("slots not passed for layout '%s", l)
		}
		if lInput.Raid == nil {
			return nil, fmt.Errorf("raid not passed for layout '%s'", l)
		}
		result = append(result, lInput)
	}
	return result, nil
}

// parsePartitions parses partitions strings into []parsedPartition
func parsePartitions(partitions []string) ([]parsedPartition, error) {
	var result []parsedPartition

	for _, p := range partitions {
		var pp parsedPartition
		parts := strings.Split(p, ",")

		for _, part := range parts {
			pair := strings.SplitN(part, "=", 2)
			if len(pair) != 2 {
				continue
			}
			key, val := pair[0], pair[1]

			switch key {
			case "slot":
				slot, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("invalid slot value: %s", val)
				}
				pp.Slots = append(pp.Slots, slot)
			case "target":
				pp.Partition.Target = val
			case "size":
				size, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("invalid size: %s", val)
				}
				pp.Partition.Size = size
			case "fs":
				pp.Partition.Fs = &val
			case "fill":
				pp.Partition.Fill = (strings.ToLower(val) == "true")
			}
		}

		if len(pp.Slots) == 0 {
			return nil, fmt.Errorf("no slot specified for partition: %s", p)
		}

		sort.Ints(pp.Slots)
		result = append(result, pp)
	}

	return result, nil
}

// mergeLayouts merges layouts by merging slots if overlaps or just appends new layout
func mergeLayouts(inputLayouts, newLayouts []serverscom.DedicatedServerLayoutInput) []serverscom.DedicatedServerLayoutInput {
	for _, newLayout := range newLayouts {
		merged := false

		for i := range inputLayouts {
			existing := &inputLayouts[i]

			slotSet := make(map[int]struct{}, len(existing.SlotPositions))
			for _, s := range existing.SlotPositions {
				slotSet[s] = struct{}{}
			}

			overlap := false
			for _, s := range newLayout.SlotPositions {
				if _, ok := slotSet[s]; ok {
					overlap = true
					break
				}
			}

			if !overlap {
				continue
			}

			for _, s := range newLayout.SlotPositions {
				if _, ok := slotSet[s]; !ok {
					existing.SlotPositions = append(existing.SlotPositions, s)
				}
			}
			sort.Ints(existing.SlotPositions)
			existing.Raid = newLayout.Raid
			merged = true
			break
		}

		if !merged {
			inputLayouts = append(inputLayouts, newLayout)
		}
	}

	return inputLayouts
}

// applyPartitions finds layout with slots matched partition slots an adds partitions to it.
// Partition will be overrided if target matched.
func applyPartitions(layouts []serverscom.DedicatedServerLayoutInput, parsed []parsedPartition) error {
	for _, p := range parsed {
		foundLayout := false

		for i := range layouts {
			lSlots := slices.Clone(layouts[i].SlotPositions)
			pSlots := slices.Clone(p.Slots)

			sort.Ints(lSlots)
			sort.Ints(pSlots)

			if !slices.Equal(lSlots, pSlots) {
				continue
			}

			foundLayout = true
			override := false

			for j := range layouts[i].Partitions {
				if layouts[i].Partitions[j].Target == p.Partition.Target {
					layouts[i].Partitions[j] = p.Partition
					override = true
					break
				}
			}

			if !override {
				layouts[i].Partitions = append(layouts[i].Partitions, p.Partition)
			}
			break
		}

		if !foundLayout {
			return fmt.Errorf("can't apply partition: no layout found with slots: %v", p.Slots)
		}
	}

	return nil
}
