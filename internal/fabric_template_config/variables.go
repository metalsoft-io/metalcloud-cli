package fabric_template_config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	fsc "github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
)

const tagSspGroup = "nvidia/ssp-group-id"

func threeTier(groups map[string][]*fsc.Device) bool {
	return len(groups["super_spine"]) > 0
}

func numericTag(dev *fsc.Device, key string) (int, bool) {
	v, ok := dev.TagsMap[key]
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, false
	}
	return n, true
}

// hgxPrefix derives the tenant hgx_subnets prefix-list match from the first
// octet of the leaf<->host pool (2-tier: <oct>.16.0.0/12, 3-tier: <oct>.0.0.0/8);
// an explicit override wins.
func hgxPrefix(config *fsc.Config, isThreeTier bool, override string) string {
	if override != "" {
		return override
	}
	pool := "172.16.0.0/12"
	if config.P2p != nil && config.P2p.Pools != nil {
		if v, ok := config.P2p.Pools["leafHost"]; ok && v != "" {
			pool = v
		}
	}
	firstOctet := byte(ipToUint(strings.SplitN(pool, "/", 2)[0]) >> 24)
	if isThreeTier {
		return fmt.Sprintf("%d.0.0.0/8", firstOctet)
	}
	return fmt.Sprintf("%d.16.0.0/12", firstOctet)
}

// leafInterfaces is a leaf's fabric-facing uplink ports (where it is interfaceA
// of a leaf<->spine link), tagged with the far-end position, in link-plan order.
func leafInterfaces(dev *fsc.Device, state *fsc.DesiredState) []map[string]interface{} {
	var out []map[string]interface{}
	for _, plan := range state.Links {
		if plan.DeviceA.Id == dev.Id {
			out = append(out, map[string]interface{}{
				"interfaceName":        plan.PortA,
				"linkedSwitchPosition": plan.DeviceB.Position,
			})
		}
	}
	return out
}

// computeFreeformVariables: every switch gets mode + hgx_prefix; an l3evpn leaf
// also gets nve_source (its loopback). Returns nil error only via ConfigError.
func computeFreeformVariables(groups map[string][]*fsc.Device, state *fsc.DesiredState, records map[int64]*deviceRecord, mode, hgx string) (map[int64]map[string]interface{}, error) {
	variables := map[int64]map[string]interface{}{}
	for _, position := range switchPositions {
		for _, dev := range groups[position] {
			vars := map[string]interface{}{"mode": mode, "hgx_prefix": hgx}
			if position == "leaf" && mode == "l3evpn" {
				loopback := loopbackOf(dev, state, records)
				if loopback == "" {
					return nil, fmt.Errorf("%s has no loopback; the l3evpn NVE source is the leaf loopback", dev.Label())
				}
				vars["nve_source"] = loopback
			}
			variables[dev.Id] = vars
		}
	}
	return variables, nil
}

// computeBgpVariables: one bgp_neighbors entry per fabric link on BOTH endpoints
// (neighbor IP = far end's /31 address), plus the leaf's per-rail /26 aggregates.
func computeBgpVariables(groups map[string][]*fsc.Device, state *fsc.DesiredState, records map[int64]*deviceRecord, mode string) (map[int64]map[string]interface{}, error) {
	isThreeTier := threeTier(groups)

	neighbors := map[int64][]map[string]interface{}{}
	for _, plan := range state.Links {
		if plan.Subnet == nil {
			return nil, fmt.Errorf("link plan has no /31 (is 'p2p' missing?); bgp needs the link subnets")
		}
		ipA, ipB := subnetHostAddrs(plan.Subnet)
		neighbors[plan.DeviceA.Id] = append(neighbors[plan.DeviceA.Id], map[string]interface{}{
			"ip": ipB, "host": hostOf(plan.DeviceB, state, records), "port": plan.PortB, "role": plan.DeviceB.Position,
		})
		neighbors[plan.DeviceB.Id] = append(neighbors[plan.DeviceB.Id], map[string]interface{}{
			"ip": ipA, "host": hostOf(plan.DeviceA, state, records), "port": plan.PortA, "role": plan.DeviceA.Position,
		})
	}

	// /26 aggregate = the leaf's host-downlink /31 run with the host-walk octet cleared.
	aggregateBases := map[int64]map[uint32]bool{}
	for _, hp := range state.HostLinks {
		if hp.Subnet == nil {
			continue
		}
		base := ipToUint(hp.Subnet.NetworkAddress) & 0xFFFFFF00
		if aggregateBases[hp.Leaf.Id] == nil {
			aggregateBases[hp.Leaf.Id] = map[uint32]bool{}
		}
		aggregateBases[hp.Leaf.Id][base] = true
	}

	variables := map[int64]map[string]interface{}{}
	for _, position := range switchPositions {
		for _, dev := range groups[position] {
			devNeighbors := append([]map[string]interface{}{}, neighbors[dev.Id]...)
			sort.SliceStable(devNeighbors, func(i, j int) bool {
				return ipToUint(devNeighbors[i]["ip"].(string)) < ipToUint(devNeighbors[j]["ip"].(string))
			})
			var bases []uint32
			for b := range aggregateBases[dev.Id] {
				bases = append(bases, b)
			}
			sort.Slice(bases, func(i, j int) bool { return bases[i] < bases[j] })
			aggregates := make([]string, 0, len(bases))
			for _, b := range bases {
				aggregates = append(aggregates, uintToIP(b)+"/26")
			}
			variables[dev.Id] = map[string]interface{}{
				"mode":          mode,
				"is_three_tier": isThreeTier,
				"aggregates":    aggregates,
				"bgp_neighbors": devNeighbors,
			}
		}
	}
	return variables, nil
}

// evpnRouteReflectors selects the EVPN overlay route reflectors (PDF 7.3.2):
// 2-tier: the 2 lowest-router-id spines; 3-tier single group: that group's 2
// lowest; 3-tier multiple groups: the lowest of each group. Router-id = loopback.
func evpnRouteReflectors(groups map[string][]*fsc.Device, state *fsc.DesiredState, records map[int64]*deviceRecord) ([]*fsc.Device, error) {
	rid := func(dev *fsc.Device) (uint32, error) {
		lb := loopbackOf(dev, state, records)
		if lb == "" {
			return 0, fmt.Errorf("%s has no loopback; overlay RR selection needs router-ids", dev.Label())
		}
		return ipToUint(lb), nil
	}
	sortByRid := func(devs []*fsc.Device) ([]*fsc.Device, error) {
		out := append([]*fsc.Device{}, devs...)
		var sortErr error
		sort.SliceStable(out, func(i, j int) bool {
			a, err1 := rid(out[i])
			b, err2 := rid(out[j])
			if err1 != nil {
				sortErr = err1
			}
			if err2 != nil {
				sortErr = err2
			}
			return a < b
		})
		return out, sortErr
	}

	ssps := groups["super_spine"]
	if len(ssps) == 0 {
		sorted, err := sortByRid(groups["spine"])
		if err != nil {
			return nil, err
		}
		if len(sorted) > 2 {
			sorted = sorted[:2]
		}
		return sorted, nil
	}

	perGroup := map[int][]*fsc.Device{}
	var groupOrder []int
	for _, dev := range ssps {
		g, ok := numericTag(dev, tagSspGroup)
		if !ok {
			return nil, fmt.Errorf("%s missing %s tag", dev.Label(), tagSspGroup)
		}
		if _, seen := perGroup[g]; !seen {
			groupOrder = append(groupOrder, g)
		}
		perGroup[g] = append(perGroup[g], dev)
	}
	if len(perGroup) == 1 {
		sorted, err := sortByRid(ssps)
		if err != nil {
			return nil, err
		}
		if len(sorted) > 2 {
			sorted = sorted[:2]
		}
		return sorted, nil
	}
	sort.Ints(groupOrder)
	var out []*fsc.Device
	for _, g := range groupOrder {
		sorted, err := sortByRid(perGroup[g])
		if err != nil {
			return nil, err
		}
		out = append(out, sorted[0])
	}
	return out, nil
}

// computeOverlayVariables: loopback-to-loopback overlay mesh. Leaves peer with
// the RR loopbacks; RRs peer with every leaf loopback; everything else gets none.
func computeOverlayVariables(groups map[string][]*fsc.Device, state *fsc.DesiredState, records map[int64]*deviceRecord, mode string) (map[int64]map[string]interface{}, error) {
	isThreeTier := threeTier(groups)
	rrs, err := evpnRouteReflectors(groups, state, records)
	if err != nil {
		return nil, err
	}
	rrIDs := map[int64]bool{}
	for _, dev := range rrs {
		rrIDs[dev.Id] = true
	}

	neighbor := func(dev *fsc.Device) (map[string]interface{}, error) {
		lb := loopbackOf(dev, state, records)
		if lb == "" {
			return nil, fmt.Errorf("%s has no loopback; overlay neighbors are loopback peerings", dev.Label())
		}
		return map[string]interface{}{"ip": lb, "host": hostOf(dev, state, records)}, nil
	}
	byIP := func(entries []map[string]interface{}) []map[string]interface{} {
		sort.SliceStable(entries, func(i, j int) bool {
			return ipToUint(entries[i]["ip"].(string)) < ipToUint(entries[j]["ip"].(string))
		})
		return entries
	}

	var rrNeighbors []map[string]interface{}
	for _, position := range []string{"spine", "super_spine"} {
		for _, dev := range groups[position] {
			if rrIDs[dev.Id] {
				n, err := neighbor(dev)
				if err != nil {
					return nil, err
				}
				rrNeighbors = append(rrNeighbors, n)
			}
		}
	}
	rrNeighbors = byIP(rrNeighbors)

	var leafNeighbors []map[string]interface{}
	for _, dev := range groups["leaf"] {
		n, err := neighbor(dev)
		if err != nil {
			return nil, err
		}
		leafNeighbors = append(leafNeighbors, n)
	}
	leafNeighbors = byIP(leafNeighbors)

	ttl := 2
	if isThreeTier {
		ttl = 3
	}
	variables := map[int64]map[string]interface{}{}
	for _, position := range switchPositions {
		for _, dev := range groups[position] {
			isRR := rrIDs[dev.Id]
			var overlay []map[string]interface{}
			switch {
			case position == "leaf":
				overlay = rrNeighbors
			case isRR:
				overlay = leafNeighbors
			}
			variables[dev.Id] = map[string]interface{}{
				"mode":                 mode,
				"is_evpn_rr":           isRR,
				"overlay_multihop_ttl": ttl,
				"overlay_neighbors":    overlay,
			}
		}
	}
	return variables, nil
}

// overlayApplies: whether the overlay template renders anything for this device.
func overlayApplies(dev *fsc.Device, vars map[string]interface{}) bool {
	return vars["mode"] == "l3evpn" && (dev.Position == "leaf" || vars["is_evpn_rr"] == true)
}

// computePfcVariables: PFC is constant lines on every switch; only var is mode.
func computePfcVariables(groups map[string][]*fsc.Device, mode string) map[int64]map[string]interface{} {
	variables := map[int64]map[string]interface{}{}
	for _, position := range switchPositions {
		for _, dev := range groups[position] {
			variables[dev.Id] = map[string]interface{}{"mode": mode}
		}
	}
	return variables
}

func pfcApplies(vars map[string]interface{}) bool { return vars["mode"] == "l3evpn" }

// renderContextFreeform / renderContextBgp build the merged variable bag the
// engine sees for a device (device fields it spreads + the profile variables),
// used by the engine-render verification.
func renderContextFreeform(dev *fsc.Device, vars map[string]interface{}, state *fsc.DesiredState, records map[int64]*deviceRecord) map[string]interface{} {
	ctx := map[string]interface{}{"position": dev.Position, "identifierString": dev.IdentifierString}
	for k, v := range vars {
		ctx[k] = v
	}
	if dev.Position == "leaf" {
		ctx["interfaces"] = leafInterfaces(dev, state)
	}
	return ctx
}

func renderContextBgp(dev *fsc.Device, vars map[string]interface{}, records map[int64]*deviceRecord) map[string]interface{} {
	ctx := map[string]interface{}{"position": dev.Position}
	if rec, ok := records[dev.Id]; ok {
		if rec.Asn != nil {
			ctx["asn"] = *rec.Asn
		}
		if rec.LoopbackAddressIpv4 != nil {
			ctx["loopbackAddress"] = *rec.LoopbackAddressIpv4
		}
	}
	for k, v := range vars {
		ctx[k] = v
	}
	return ctx
}
