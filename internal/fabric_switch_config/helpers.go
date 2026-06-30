package fabric_switch_config

import (
	"sort"
	"strings"
)

func (s *DesiredState) warn(msg string) { s.Warnings = append(s.Warnings, msg) }

// sortByNumericTags returns a stably-sorted copy of devs ordered by the numeric
// tag tuple (keys). Equal tuples keep the input (ordering) order.
func sortByNumericTags(devs []*Device, keys []string) ([]*Device, error) {
	type kv struct {
		dev *Device
		key []int
	}
	arr := make([]kv, len(devs))
	for i, d := range devs {
		k, err := numericTags(d, keys...)
		if err != nil {
			return nil, err
		}
		arr[i] = kv{d, k}
	}
	sort.SliceStable(arr, func(i, j int) bool { return compareIntTuples(arr[i].key, arr[j].key) })
	out := make([]*Device, len(arr))
	for i := range arr {
		out[i] = arr[i].dev
	}
	return out, nil
}

type pairKey struct {
	a, b int
	s    string
}

// groupByTagPair groups devices by the (keyA, keyB) numeric tag pair, returning
// the groups (keyed by a string signature) and the distinct keys sorted by
// (a, b).
func groupByTagPair(devs []*Device, keyA, keyB string) (map[string][]*Device, []pairKey, error) {
	groups := map[string][]*Device{}
	seen := map[string]pairKey{}
	for _, d := range devs {
		vals, err := numericTags(d, keyA, keyB)
		if err != nil {
			return nil, nil, err
		}
		pk := pairKey{a: vals[0], b: vals[1]}
		pk.s = sigOf(vals)
		groups[pk.s] = append(groups[pk.s], d)
		seen[pk.s] = pk
	}
	order := make([]pairKey, 0, len(seen))
	for _, pk := range seen {
		order = append(order, pk)
	}
	sort.SliceStable(order, func(i, j int) bool {
		if order[i].a != order[j].a {
			return order[i].a < order[j].a
		}
		return order[i].b < order[j].b
	})
	return groups, order, nil
}

func sigOf(vals []int) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = itoa(int64(v))
	}
	return strings.Join(parts, "|")
}

func equalIntSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func allEqual(xs []int) bool {
	for i := 1; i < len(xs); i++ {
		if xs[i] != xs[0] {
			return false
		}
	}
	return true
}

func intsOf(devs []*Device, key string) []int {
	out := make([]int, 0, len(devs))
	for _, d := range devs {
		n, _ := numericTag(d, key)
		out = append(out, n)
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedKeysInt(m map[int][]*Device) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}

func keySet(m map[int][]*Device) map[int]bool {
	out := map[int]bool{}
	for k := range m {
		out[k] = true
	}
	return out
}

func equalIntSets(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func sortedSet(m map[int]bool) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}

func joinStrings(xs []string, sep string) string { return strings.Join(xs, sep) }
