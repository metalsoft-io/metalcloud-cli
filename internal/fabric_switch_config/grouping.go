package fabric_switch_config

import (
	"sort"
	"strconv"
	"strings"
)

// Valid ordering keys (the stable order that defines each device's ordinal
// within its position group).
const (
	OrderingManagementAddress = "managementAddress"
	OrderingIdentifierString  = "identifierString"
	OrderingId                = "id"
)

var validOrderings = []string{OrderingManagementAddress, OrderingIdentifierString, OrderingId}

// GroupAndOrder groups devices by position and sorts each group by the
// configured ordering key. A device's 1-based ordinal is its index+1 in its
// group's slice.
func GroupAndOrder(devices []*Device, ordering string) (map[string][]*Device, error) {
	groups := map[string][]*Device{}
	for _, dev := range devices {
		position := dev.Position
		if position == "" {
			position = "unknown"
		}
		groups[position] = append(groups[position], dev)
	}

	var sortErr error
	for _, devs := range groups {
		// Validate up front so sort comparisons stay total.
		if ordering == OrderingManagementAddress {
			for _, dev := range devs {
				if _, ok := ipv4ToUint(dev.ManagementAddress); !ok {
					return nil, configErrorf(
						"device id=%d (%q) has no managementAddress; use 'ordering: id' or fix the device",
						dev.Id, dev.IdentifierString)
				}
			}
		}
		devs := devs
		sort.SliceStable(devs, func(i, j int) bool {
			switch ordering {
			case OrderingManagementAddress:
				a, _ := ipv4ToUint(devs[i].ManagementAddress)
				b, _ := ipv4ToUint(devs[j].ManagementAddress)
				return a < b
			case OrderingIdentifierString:
				return devs[i].IdentifierString < devs[j].IdentifierString
			default: // id
				return devs[i].Id < devs[j].Id
			}
		})
	}
	if sortErr != nil {
		return nil, sortErr
	}
	return groups, nil
}

func isValidOrdering(ordering string) bool {
	for _, o := range validOrderings {
		if o == ordering {
			return true
		}
	}
	return false
}

// numericTag reads a required numeric tag off the device; missing or
// non-numeric is fatal.
func numericTag(dev *Device, key string) (int, error) {
	value, ok := dev.TagsMap[key]
	if !ok {
		return 0, configErrorf("tag '%s' not present in tagsMap of %s", key, dev.Label())
	}
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, configErrorf("tag '%s' of %s is not numeric: %q", key, dev.Label(), value)
	}
	return n, nil
}

// numericTags reads several numeric tags and returns them as a key tuple.
func numericTags(dev *Device, keys ...string) ([]int, error) {
	out := make([]int, len(keys))
	for i, k := range keys {
		n, err := numericTag(dev, k)
		if err != nil {
			return nil, err
		}
		out[i] = n
	}
	return out, nil
}

// compareIntTuples returns true if a < b lexicographically.
func compareIntTuples(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}
