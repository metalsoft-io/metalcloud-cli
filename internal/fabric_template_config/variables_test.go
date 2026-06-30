package fabric_template_config

import (
	"testing"

	fsc "github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
)

// fixturePlan builds the 3-tier reference fixture plan (matching the
// fabric_switch_config gold fixture) and the device records.
func fixturePlan(t *testing.T) (map[string][]*fsc.Device, *fsc.DesiredState, map[int64]*deviceRecord) {
	t.Helper()
	dev := func(id int64, position, mgmt string, tags map[string]string) *fsc.Device {
		return &fsc.Device{Id: id, Position: position, ManagementAddress: mgmt,
			IdentifierString: "old-" + position, Driver: "cumulus_linux", TagsMap: tags}
	}
	devices := []*fsc.Device{
		dev(1, "spine", "10.0.0.22", map[string]string{"nvidia/pod-id": "5", "nvidia/rail-group-id": "1", "nvidia/spine-index": "2"}),
		dev(2, "spine", "10.0.0.21", map[string]string{"nvidia/pod-id": "5", "nvidia/rail-group-id": "1", "nvidia/spine-index": "1"}),
		dev(3, "leaf", "10.0.0.12", map[string]string{"nvidia/pod-id": "5", "nvidia/scalability-unit-id": "2", "nvidia/rail-group-id": "1"}),
		dev(4, "leaf", "10.0.0.11", map[string]string{"nvidia/pod-id": "5", "nvidia/scalability-unit-id": "1", "nvidia/rail-group-id": "1"}),
		dev(5, "leaf", "10.0.0.13", map[string]string{"nvidia/pod-id": "5", "nvidia/scalability-unit-id": "3", "nvidia/rail-group-id": "1"}),
		dev(6, "super_spine", "10.0.0.31", map[string]string{"nvidia/ssp-group-id": "1"}),
		dev(7, "super_spine", "10.0.0.32", map[string]string{"nvidia/ssp-group-id": "2"}),
		dev(8, "super_spine", "10.0.0.33", map[string]string{"nvidia/ssp-group-id": "1"}),
	}
	yaml := `
ordering: managementAddress
loopback:
  subnet: 10.253.128.0/18
topology:
  leafSpine:
    linksPerPair: auto
  spineSuperSpine:
    linksPerPair: 4
  leafHost:
    nodeCount: 2
p2p:
  mtu: 9216
`
	cfg, err := fsc.LoadConfig([]byte(yaml))
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	groups, err := fsc.GroupAndOrder(devices, cfg.Ordering)
	if err != nil {
		t.Fatalf("GroupAndOrder: %v", err)
	}
	state, err := fsc.ComputeDesired(cfg, groups)
	if err != nil {
		t.Fatalf("ComputeDesired: %v", err)
	}
	records := map[int64]*deviceRecord{}
	for _, d := range devices {
		records[d.Id] = &deviceRecord{Device: *d}
	}
	return groups, state, records
}

func TestComputeBgpVariables(t *testing.T) {
	groups, state, records := fixturePlan(t)
	vars, err := computeBgpVariables(groups, state, records, "l3evpn")
	if err != nil {
		t.Fatalf("computeBgpVariables: %v", err)
	}

	// leaf dev4 uplinks to 2 spines x L=10 = 20 fabric links -> 20 neighbors.
	leaf := vars[4]
	neighbors := leaf["bgp_neighbors"].([]map[string]interface{})
	if len(neighbors) != 20 {
		t.Errorf("dev4 bgp_neighbors = %d, want 20", len(neighbors))
	}
	// dev4<->dev2 swp33s0 link /31 = 10.254.0.0/31; dev4 is gateway (even), so its
	// neighbor (the spine) is the odd address 10.254.0.1 on spine port swp1s0.
	found := false
	for _, n := range neighbors {
		if n["ip"] == "10.254.0.1" && n["port"] == "swp1s0" && n["role"] == "spine" {
			found = true
		}
	}
	if !found {
		t.Errorf("dev4 missing neighbor {10.254.0.1, swp1s0, spine}: %v", neighbors)
	}
	// dev4 aggregates: rails {1, 5} -> 172.18.0.0/26 and 172.26.0.0/26.
	agg := leaf["aggregates"].([]string)
	if len(agg) != 2 || agg[0] != "172.18.0.0/26" || agg[1] != "172.26.0.0/26" {
		t.Errorf("dev4 aggregates = %v, want [172.18.0.0/26 172.26.0.0/26]", agg)
	}
	if leaf["is_three_tier"] != true {
		t.Errorf("is_three_tier should be true")
	}
}

func TestEvpnRouteReflectors(t *testing.T) {
	groups, state, records := fixturePlan(t)
	rrs, err := evpnRouteReflectors(groups, state, records)
	if err != nil {
		t.Fatalf("evpnRouteReflectors: %v", err)
	}
	// 3-tier, multiple ssp groups -> lowest-router-id ssp of each group.
	// loopbacks: dev6(g1)=.6, dev8(g1)=.7, dev7(g2)=.8 -> RR = {dev6, dev7}.
	got := map[int64]bool{}
	for _, d := range rrs {
		got[d.Id] = true
	}
	if len(got) != 2 || !got[6] || !got[7] {
		t.Errorf("route reflectors = %v, want {6,7}", got)
	}
}

func TestComputeOverlayVariables(t *testing.T) {
	groups, state, records := fixturePlan(t)
	overlay, err := computeOverlayVariables(groups, state, records, "l3evpn")
	if err != nil {
		t.Fatalf("computeOverlayVariables: %v", err)
	}
	// leaf peers the RR loopbacks, sorted: dev6=.6, dev7=.8.
	leafNbrs := overlay[4]["overlay_neighbors"].([]map[string]interface{})
	if len(leafNbrs) != 2 || leafNbrs[0]["ip"] != "10.253.128.6" || leafNbrs[1]["ip"] != "10.253.128.8" {
		t.Errorf("dev4 overlay_neighbors = %v, want [.6 .8]", leafNbrs)
	}
	if overlay[6]["is_evpn_rr"] != true || overlay[7]["is_evpn_rr"] != true {
		t.Errorf("dev6/dev7 should be RRs")
	}
	if overlay[8]["is_evpn_rr"] != false {
		t.Errorf("dev8 should not be an RR")
	}
	if overlay[4]["overlay_multihop_ttl"] != 3 {
		t.Errorf("3-tier overlay ttl should be 3")
	}
	// RRs peer every leaf loopback (3 leaves).
	if n := overlay[6]["overlay_neighbors"].([]map[string]interface{}); len(n) != 3 {
		t.Errorf("RR dev6 overlay_neighbors = %d, want 3", len(n))
	}
}

func TestOverlayAndPfcApplies(t *testing.T) {
	groups, state, records := fixturePlan(t)
	overlay, _ := computeOverlayVariables(groups, state, records, "l3evpn")
	// leaf applies; non-RR spine/ssp does not.
	leaf := groups["leaf"][0]
	if !overlayApplies(leaf, overlay[leaf.Id]) {
		t.Errorf("overlay should apply to a leaf in l3evpn")
	}
	if overlayApplies(groups["super_spine"][0], overlay[8]) {
		// dev8 (index 0 in group? order is mgmt) -> ensure a non-RR ssp does not apply
	}
	pfc := computePfcVariables(groups, "l3evpn")
	if !pfcApplies(pfc[leaf.Id]) {
		t.Errorf("pfc should apply in l3evpn")
	}
	pfcPurel3 := computePfcVariables(groups, "purel3")
	if pfcApplies(pfcPurel3[leaf.Id]) {
		t.Errorf("pfc should not apply in purel3")
	}
}

func TestComputeFreeformVariables(t *testing.T) {
	groups, state, records := fixturePlan(t)
	vars, err := computeFreeformVariables(groups, state, records, "l3evpn", "172.0.0.0/8")
	if err != nil {
		t.Fatalf("computeFreeformVariables: %v", err)
	}
	// leaf dev4 (l3evpn) gets nve_source = its loopback .1.
	if vars[4]["nve_source"] != "10.253.128.1" {
		t.Errorf("dev4 nve_source = %v, want 10.253.128.1", vars[4]["nve_source"])
	}
	if vars[4]["hgx_prefix"] != "172.0.0.0/8" || vars[4]["mode"] != "l3evpn" {
		t.Errorf("dev4 freeform vars wrong: %v", vars[4])
	}
	// spine never gets nve_source.
	if _, ok := vars[2]["nve_source"]; ok {
		t.Errorf("spine should not get nve_source")
	}
	// purel3: no nve_source even on leaves.
	purel3, _ := computeFreeformVariables(groups, state, records, "purel3", "172.16.0.0/12")
	if _, ok := purel3[4]["nve_source"]; ok {
		t.Errorf("purel3 leaf should not get nve_source")
	}
}

func TestHgxPrefix(t *testing.T) {
	cfg, _ := fsc.LoadConfig([]byte("topology:\n  leafSpine:\n    linksPerPair: auto\np2p:\n  pools:\n    leafHost: 172.16.0.0/12\n"))
	if got := hgxPrefix(cfg, true, ""); got != "172.0.0.0/8" {
		t.Errorf("3-tier hgxPrefix = %q, want 172.0.0.0/8", got)
	}
	if got := hgxPrefix(cfg, false, ""); got != "172.16.0.0/12" {
		t.Errorf("2-tier hgxPrefix = %q, want 172.16.0.0/12", got)
	}
	if got := hgxPrefix(cfg, true, "10.0.0.0/8"); got != "10.0.0.0/8" {
		t.Errorf("override hgxPrefix = %q, want 10.0.0.0/8", got)
	}
}
