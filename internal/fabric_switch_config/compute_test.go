package fabric_switch_config

import (
	"testing"
)

func ptrInt(v int) *int       { return &v }
func ptrInt32(v int32) *int32 { return &v }
func ptrStr(v string) *string { return &v }

func fixtureDevices() []*Device {
	d := func(id int64, position, mgmt string, tags map[string]string) *Device {
		return &Device{
			Id: id, Position: position, ManagementAddress: mgmt,
			IdentifierString: "old-" + position, Driver: "cumulus_linux", TagsMap: tags,
		}
	}
	return []*Device{
		d(1, "spine", "10.0.0.22", map[string]string{tagPod: "5", tagRail: "1", tagSpineIndex: "2"}),
		d(2, "spine", "10.0.0.21", map[string]string{tagPod: "5", tagRail: "1", tagSpineIndex: "1"}),
		d(3, "leaf", "10.0.0.12", map[string]string{tagPod: "5", tagSu: "2", tagRail: "1"}),
		d(4, "leaf", "10.0.0.11", map[string]string{tagPod: "5", tagSu: "1", tagRail: "1"}),
		d(5, "leaf", "10.0.0.13", map[string]string{tagPod: "5", tagSu: "3", tagRail: "1"}),
		d(6, "super_spine", "10.0.0.31", map[string]string{tagSspGroup: "1"}),
		d(7, "super_spine", "10.0.0.32", map[string]string{tagSspGroup: "2"}),
		d(8, "super_spine", "10.0.0.33", map[string]string{tagSspGroup: "1"}),
	}
}

func fixtureConfig() *Config {
	return &Config{
		Ordering: OrderingManagementAddress,
		Hostname: &HostnameConfig{Templates: map[string]*string{
			"leaf":        ptrStr("leaf-pod{tag:nvidia/pod-id}-su{tag:nvidia/scalability-unit-id}-r{tag:nvidia/rail-group-id}"),
			"spine":       ptrStr("spine-pod{tag:nvidia/pod-id}-r{tag:nvidia/rail-group-id}-s{tag:nvidia/spine-index}"),
			"super_spine": ptrStr("ssp-group{tag:nvidia/ssp-group-id}-s{ordinalBy:nvidia/ssp-group-id}"),
		}},
		Asn:      &AsnConfig{},
		Loopback: &LoopbackConfig{},
		Topology: &TopologyConfig{
			LeafSpine:       &LayerConfig{},
			SpineSuperSpine: &LayerConfig{LinksPerPair: ptrInt(4)},
			LeafHost:        &LeafHostConfig{NodeCount: ptrInt(2)},
		},
		P2p:                 &P2pConfig{Mtu: ptrInt32(9216)},
		DescriptionTemplate: ptrStr("to-{peerHostname}:{peerPort}"),
	}
}

func computeFixture(t *testing.T) *DesiredState {
	t.Helper()
	groups, err := GroupAndOrder(fixtureDevices(), OrderingManagementAddress)
	if err != nil {
		t.Fatalf("GroupAndOrder: %v", err)
	}
	state, err := ComputeDesired(fixtureConfig(), groups)
	if err != nil {
		t.Fatalf("ComputeDesired: %v", err)
	}
	return state
}

func hostname(t *testing.T, s *DesiredState, id int64) string {
	t.Helper()
	d := s.ByDevice[id]
	if d == nil || d.Hostname == nil {
		t.Fatalf("device %d has no hostname", id)
	}
	return *d.Hostname
}

func TestHostnames(t *testing.T) {
	s := computeFixture(t)
	want := map[int64]string{
		2: "spine-pod5-r1-s1", 1: "spine-pod5-r1-s2",
		4: "leaf-pod5-su1-r1", 3: "leaf-pod5-su2-r1", 5: "leaf-pod5-su3-r1",
		6: "ssp-group1-s1", 7: "ssp-group2-s1", 8: "ssp-group1-s2",
	}
	for id, w := range want {
		if got := hostname(t, s, id); got != w {
			t.Errorf("device %d hostname = %q, want %q", id, got, w)
		}
	}
}

func TestAsns(t *testing.T) {
	s := computeFixture(t)
	want := map[int64]int64{
		4: 4200000000, 3: 4200000001, 5: 4200000002, // leaves unique by (pod,su,rail)
		1: 4201000000, 2: 4201000000, // spines share per (pod,rail)
		6: 4202000000, 7: 4202000000, 8: 4202000000, // ssp all equal
	}
	for id, w := range want {
		d := s.ByDevice[id]
		if d == nil || d.Asn == nil {
			t.Fatalf("device %d has no asn", id)
		}
		if *d.Asn != w {
			t.Errorf("device %d asn = %d, want %d", id, *d.Asn, w)
		}
	}
}

func TestLoopbacks(t *testing.T) {
	s := computeFixture(t)
	// leaves r1,r2,r3 then spines s1,s2 then ssp g1#1,g1#2,g2#1
	want := []struct {
		id int64
		ip string
	}{
		{4, "10.253.128.1"}, {3, "10.253.128.2"}, {5, "10.253.128.3"},
		{2, "10.253.128.4"}, {1, "10.253.128.5"},
		{6, "10.253.128.6"}, {8, "10.253.128.7"}, {7, "10.253.128.8"},
	}
	for _, w := range want {
		d := s.ByDevice[w.id]
		if d == nil || d.LoopbackIp == nil {
			t.Fatalf("device %d has no loopback", w.id)
		}
		if *d.LoopbackIp != w.ip {
			t.Errorf("device %d loopback = %q, want %q", w.id, *d.LoopbackIp, w.ip)
		}
	}
}

func findLink(s *DesiredState, layer string, aId, bId int64, portA string) *LinkPlan {
	for _, l := range s.Links {
		if l.Layer == layer && l.DeviceA.Id == aId && l.DeviceB.Id == bId && l.PortA == portA {
			return l
		}
	}
	return nil
}

func TestLinkCounts(t *testing.T) {
	s := computeFixture(t)
	ls, ssp := 0, 0
	for _, l := range s.Links {
		switch l.Layer {
		case "leafSpine":
			ls++
		case "spineSuperSpine":
			ssp++
		}
	}
	if ls != 60 || ssp != 12 {
		t.Errorf("link counts = (%d leafSpine, %d spineSuperSpine), want (60, 12)", ls, ssp)
	}
}

func TestLeafSpineLinks(t *testing.T) {
	s := computeFixture(t)
	cases := []struct {
		aId, bId             int64
		portA, portB, subnet string
	}{
		{4, 2, "swp33s0", "swp1s0", "10.254.0.0/31"},   // block0<->global0 u0
		{3, 2, "swp37s1", "swp10s1", "10.254.0.38/31"}, // leaf block offset, u=L-1
		{5, 1, "swp38s0", "swp11s0", "10.254.1.40/31"}, // spine-global third octet
	}
	for _, c := range cases {
		l := findLink(s, "leafSpine", c.aId, c.bId, c.portA)
		if l == nil {
			t.Errorf("leafSpine link %d->%d %s not found", c.aId, c.bId, c.portA)
			continue
		}
		if l.PortB != c.portB || l.Subnet.String() != c.subnet {
			t.Errorf("leafSpine %d->%d %s = (%s, %s), want (%s, %s)",
				c.aId, c.bId, c.portA, l.PortB, l.Subnet, c.portB, c.subnet)
		}
	}
}

func TestSpineSuperSpineLinks(t *testing.T) {
	s := computeFixture(t)
	cases := []struct {
		aId, bId             int64
		portA, portB, subnet string
	}{
		{2, 6, "swp17s0", "swp1s0", "100.64.0.0/31"},  // in-group block 0
		{2, 8, "swp20s1", "swp2s1", "100.64.0.14/31"}, // in-group block 1, u=3
		{1, 7, "swp17s0", "swp3s0", "100.64.0.64/31"}, // ssp port by global spine index + 64-run
	}
	for _, c := range cases {
		l := findLink(s, "spineSuperSpine", c.aId, c.bId, c.portA)
		if l == nil {
			t.Errorf("spineSuperSpine link %d->%d %s not found", c.aId, c.bId, c.portA)
			continue
		}
		if l.PortB != c.portB || l.Subnet.String() != c.subnet {
			t.Errorf("spineSuperSpine %d->%d %s = (%s, %s), want (%s, %s)",
				c.aId, c.bId, c.portA, l.PortB, l.Subnet, c.portB, c.subnet)
		}
	}
}

func TestSubnetUniquenessAndContainment(t *testing.T) {
	s := computeFixture(t)
	seen := map[string]bool{}
	for _, l := range s.Links {
		seen[l.Subnet.String()] = true
	}
	for _, h := range s.HostLinks {
		seen[h.Subnet.String()] = true
	}
	if len(seen) != 60+12+12 {
		t.Errorf("unique /31 count = %d, want %d", len(seen), 84)
	}
	leafSpinePool, _ := parseIpv4Network("10.254.0.0/16")
	sspPool, _ := parseIpv4Network("100.64.0.0/10")
	hostPool, _ := parseIpv4Network("172.16.0.0/12")
	for _, l := range s.Links {
		base, _ := ipv4ToUint(l.Subnet.NetworkAddress)
		pool := leafSpinePool
		if l.Layer == "spineSuperSpine" {
			pool = sspPool
		}
		if !pool.containsSubnet(base, 31) {
			t.Errorf("%s /31 %s outside its pool", l.Layer, l.Subnet)
		}
	}
	for _, h := range s.HostLinks {
		base, _ := ipv4ToUint(h.Subnet.NetworkAddress)
		if !hostPool.containsSubnet(base, 31) {
			t.Errorf("host /31 %s outside leafHost pool", h.Subnet)
		}
	}
}

func hostByPort(s *DesiredState, leafId int64, port string) *HostLinkPlan {
	for _, h := range s.HostLinks {
		if h.Leaf.Id == leafId && h.LeafPort == port {
			return h
		}
	}
	return nil
}

func TestHostLinkFormula(t *testing.T) {
	s := computeFixture(t)
	cases := []struct {
		leafId int64
		port   string
		subnet string
	}{
		{4, "swp1s0", "172.18.0.0/31"},
		{4, "swp1s1", "172.26.0.0/31"},
		{4, "swp2s0", "172.18.0.2/31"},
		{5, "swp1s0", "172.18.2.0/31"},
	}
	for _, c := range cases {
		h := hostByPort(s, c.leafId, c.port)
		if h == nil {
			t.Errorf("host link %d %s not found", c.leafId, c.port)
			continue
		}
		if h.Subnet.String() != c.subnet {
			t.Errorf("host /31 %d %s = %s, want %s", c.leafId, c.port, h.Subnet, c.subnet)
		}
	}
}

func TestHostLinkMetadata(t *testing.T) {
	s := computeFixture(t)
	h := hostByPort(s, 4, "swp1s0")
	if h.HostName != "hgx-pod05-su01-h00" || h.Nic != "enp60s0f0np0" {
		t.Errorf("host meta = (%q, %q), want (hgx-pod05-su01-h00, enp60s0f0np0)", h.HostName, h.Nic)
	}
	if h1 := hostByPort(s, 4, "swp1s1"); h1.Nic != "enp188s0f0np0" {
		t.Errorf("s1 nic = %q, want enp188s0f0np0", h1.Nic)
	}
}

func TestDescriptions(t *testing.T) {
	s := computeFixture(t)
	want := map[PortKey]string{
		{4, "swp33s0"}: "to-spine-pod5-r1-s1:swp1s0",
		{2, "swp1s0"}:  "to-leaf-pod5-su1-r1:swp33s0",
		{2, "swp17s0"}: "to-ssp-group1-s1:swp1s0",
		{6, "swp1s0"}:  "to-spine-pod5-r1-s1:swp17s0",
		{7, "swp3s0"}:  "to-spine-pod5-r1-s2:swp17s0",
		// leaf->host (3-tier default = pod-qualified)
		{4, "swp1s0"}: "to_hgx-pod05-su01-h00_enp60s0f0np0",
		{4, "swp1s1"}: "to_hgx-pod05-su01-h00_enp188s0f0np0",
		{4, "swp2s0"}: "to_hgx-pod05-su01-h01_enp60s0f0np0",
		{5, "swp1s0"}: "to_hgx-pod05-su03-h00_enp60s0f0np0",
	}
	for k, w := range want {
		if got := s.PortDescriptions[k]; got != w {
			t.Errorf("description %v = %q, want %q", k, got, w)
		}
	}
}
