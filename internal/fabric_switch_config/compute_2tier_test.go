package fabric_switch_config

import "testing"

// twoTierDevices mirrors make_2tier_fixture() in test_offline.py.
func twoTierDevices() []*Device {
	d := func(id int64, position, mgmt string, tags map[string]string) *Device {
		return &Device{Id: id, Position: position, ManagementAddress: mgmt,
			IdentifierString: "old-" + position, Driver: "cumulus_linux", TagsMap: tags}
	}
	return []*Device{
		d(21, "spine", "10.1.0.22", map[string]string{tagSpineIndex: "1"}),
		d(22, "spine", "10.1.0.21", map[string]string{tagSpineIndex: "0"}),
		d(23, "leaf", "10.1.0.12", map[string]string{tagSu: "0", tagRail: "1"}),
		d(24, "leaf", "10.1.0.11", map[string]string{tagSu: "0", tagRail: "0"}),
		d(25, "leaf", "10.1.0.14", map[string]string{tagSu: "0", tagRail: "3"}),
		d(26, "leaf", "10.1.0.13", map[string]string{tagSu: "0", tagRail: "2"}),
	}
}

func compute(t *testing.T, devices []*Device, config *Config) *DesiredState {
	t.Helper()
	groups, err := GroupAndOrder(devices, config.ordering())
	if err != nil {
		t.Fatalf("GroupAndOrder: %v", err)
	}
	state, err := ComputeDesired(config, groups)
	if err != nil {
		t.Fatalf("ComputeDesired: %v", err)
	}
	return state
}

func TestThreeTierDefaultHostnames(t *testing.T) {
	s := compute(t, fixtureDevices(), &Config{Hostname: &HostnameConfig{Templates: map[string]*string{}}})
	want := map[int64]string{
		4: "leaf-pod05-su01-r1", 2: "spine-pod05-r1-s01", 1: "spine-pod05-r1-s02",
		6: "ssp-group01-s00", 7: "ssp-group02-s00", 8: "ssp-group01-s01",
	}
	for id, w := range want {
		if got := hostname(t, s, id); got != w {
			t.Errorf("device %d default hostname = %q, want %q", id, got, w)
		}
	}
}

func TestTwoTierNamingAsnLoopback(t *testing.T) {
	cfg := &Config{
		Ordering: OrderingManagementAddress,
		Hostname: &HostnameConfig{Templates: map[string]*string{}},
		Asn:      &AsnConfig{},
		Loopback: &LoopbackConfig{},
		Topology: &TopologyConfig{LeafHost: &LeafHostConfig{NodeCount: ptrInt(1)}},
	}
	s := compute(t, twoTierDevices(), cfg)

	hostnames := map[int64]string{
		24: "leaf-su00-r0", 23: "leaf-su00-r1", 26: "leaf-su00-r2", 25: "leaf-su00-r3",
		22: "spine-s00", 21: "spine-s01",
	}
	for id, w := range hostnames {
		if got := hostname(t, s, id); got != w {
			t.Errorf("2-tier hostname %d = %q, want %q", id, got, w)
		}
	}
	asns := map[int64]int64{24: 4200000000, 23: 4200000001, 26: 4200000002, 25: 4200000003, 22: 4201000000, 21: 4201000000}
	for id, w := range asns {
		if *s.ByDevice[id].Asn != w {
			t.Errorf("2-tier asn %d = %d, want %d", id, *s.ByDevice[id].Asn, w)
		}
	}
	loopbacks := []struct {
		id int64
		ip string
	}{{24, "10.253.128.1"}, {23, "10.253.128.2"}, {26, "10.253.128.3"}, {25, "10.253.128.4"}, {22, "10.253.128.5"}, {21, "10.253.128.6"}}
	for _, w := range loopbacks {
		if *s.ByDevice[w.id].LoopbackIp != w.ip {
			t.Errorf("2-tier loopback %d = %q, want %q", w.id, *s.ByDevice[w.id].LoopbackIp, w.ip)
		}
	}
	// 2-tier host descriptions stay flat (no pod part).
	if got := s.PortDescriptions[PortKey{24, "swp1s0"}]; got != "to_hgx-su00-h00_enp26s0f0np0" {
		t.Errorf("2-tier host desc (24,swp1s0) = %q", got)
	}
	if got := s.PortDescriptions[PortKey{23, "swp1s1"}]; got != "to_hgx-su00-h00_enp188s0f0np0" {
		t.Errorf("2-tier host desc (23,swp1s1) = %q", got)
	}
}

func TestHostnameOverrideAndNullSkip(t *testing.T) {
	cfg := &Config{
		Ordering: OrderingManagementAddress,
		Hostname: &HostnameConfig{Templates: map[string]*string{
			"leaf":  ptrStr("custom-{ordinal}"),
			"spine": nil, // explicit null -> skip
		}},
	}
	s := compute(t, twoTierDevices(), cfg)
	if got := hostname(t, s, 24); got != "custom-1" {
		t.Errorf("override hostname (24) = %q, want custom-1", got)
	}
	if d := s.ByDevice[22]; d != nil && d.Hostname != nil {
		t.Errorf("spine 22 should be skipped, got %q", *d.Hostname)
	}
	if d := s.ByDevice[21]; d != nil && d.Hostname != nil {
		t.Errorf("spine 21 should be skipped, got %q", *d.Hostname)
	}
}

func TestTwoTierFullMesh(t *testing.T) {
	cfg := &Config{
		Ordering:            OrderingManagementAddress,
		Hostname:            &HostnameConfig{Templates: map[string]*string{}},
		Topology:            &TopologyConfig{LeafSpine: &LayerConfig{}, LeafHost: &LeafHostConfig{NodeCount: ptrInt(9)}},
		P2p:                 &P2pConfig{Mtu: ptrInt32(9216)},
		DescriptionTemplate: ptrStr("to_{peerHostname}_{peerPort}"),
	}
	s := compute(t, twoTierDevices(), cfg)

	mesh := 0
	for _, l := range s.Links {
		if l.Layer == "leafSpine" {
			mesh++
		}
	}
	if mesh != 256 { // 4 leaves x 2 spines x L=128//4=32
		t.Errorf("full-mesh leafSpine links = %d, want 256", mesh)
	}

	cases := []struct {
		aId, bId             int64
		portA, portB, subnet string
	}{
		{23, 22, "swp33s0", "swp17s0", "10.254.0.64/31"},  // section 1 u=0
		{23, 22, "swp48s1", "swp32s1", "10.254.0.126/31"}, // section 1 u=31
		{23, 21, "swp49s0", "swp17s0", "10.254.1.64/31"},  // spine-s01 u=0
	}
	for _, c := range cases {
		l := findLink(s, "leafSpine", c.aId, c.bId, c.portA)
		if l == nil {
			t.Errorf("mesh link %d->%d %s not found", c.aId, c.bId, c.portA)
			continue
		}
		if l.PortB != c.portB || l.Subnet.String() != c.subnet {
			t.Errorf("mesh %d->%d %s = (%s,%s), want (%s,%s)", c.aId, c.bId, c.portA, l.PortB, l.Subnet, c.portB, c.subnet)
		}
	}

	if got := s.PortDescriptions[PortKey{23, "swp33s0"}]; got != "to_spine-s00_swp17s0" {
		t.Errorf("mesh desc (23,swp33s0) = %q", got)
	}
	if got := s.PortDescriptions[PortKey{22, "swp17s0"}]; got != "to_leaf-su00-r1_swp33s0" {
		t.Errorf("mesh desc (22,swp17s0) = %q", got)
	}

	// section 3 host worked examples + flat names.
	host := map[string]string{}
	for _, h := range s.HostLinks {
		if h.Leaf.Id == 24 {
			host[h.LeafPort] = h.Subnet.String()
		}
	}
	for port, want := range map[string]string{"swp1s0": "172.16.0.0/31", "swp1s1": "172.24.0.0/31", "swp9s0": "172.16.0.16/31"} {
		if host[port] != want {
			t.Errorf("mesh host /31 (24,%s) = %q, want %q", port, host[port], want)
		}
	}
	if h := hostByPort(s, 24, "swp9s0"); h.HostName != "hgx-su00-h08" {
		t.Errorf("mesh host name (24,swp9s0) = %q, want hgx-su00-h08", h.HostName)
	}
}

func TestExplicitNodesList(t *testing.T) {
	cfg := &Config{
		Ordering: OrderingManagementAddress,
		Hostname: &HostnameConfig{Templates: map[string]*string{}},
		Topology: &TopologyConfig{LeafSpine: &LayerConfig{}, LeafHost: &LeafHostConfig{Nodes: []int{0, 8}}},
		P2p:      &P2pConfig{Mtu: ptrInt32(9216)},
	}
	s := compute(t, twoTierDevices(), cfg)
	if len(s.HostLinks) != 4*4 { // 4 leaves x 2 nodes x 2 sub-ports
		t.Errorf("host links with explicit nodes = %d, want 16", len(s.HostLinks))
	}
	if h := hostByPort(s, 24, "swp9s0"); h == nil || h.Subnet.String() != "172.16.0.16/31" || h.HostName != "hgx-su00-h08" {
		t.Errorf("explicit-nodes /31 offset not preserved for (24,swp9s0): %+v", h)
	}
}

func TestErrorPaths(t *testing.T) {
	mustErr := func(name string, devices []*Device, cfg *Config, wantSub string) {
		groups, err := GroupAndOrder(devices, cfg.ordering())
		if err == nil {
			_, err = ComputeDesired(cfg, groups)
		}
		if err == nil {
			t.Errorf("%s: expected error, got nil", name)
			return
		}
		if wantSub != "" && !contains(err.Error(), wantSub) {
			t.Errorf("%s: error %q does not contain %q", name, err.Error(), wantSub)
		}
	}

	// spineSuperSpine on a 2-tier fabric.
	mustErr("ssp on 2-tier", twoTierDevices(), &Config{
		Topology: &TopologyConfig{SpineSuperSpine: &LayerConfig{}},
	}, "super_spine")

	// missing required tag for ASN sort.
	noSu := twoTierDevices()
	delete(noSu[2].TagsMap, tagSu)
	mustErr("missing su tag", noSu, &Config{Asn: &AsnConfig{}}, "scalability-unit-id")

	// duplicate spine (pod,rail,spine-index) in 3-tier topology.
	dupSpine := fixtureDevices()
	dupSpine[0].TagsMap[tagSpineIndex] = "1" // dev1 now collides with dev2
	mustErr("duplicate spine key", dupSpine, &Config{
		Topology: &TopologyConfig{LeafSpine: &LayerConfig{}},
	}, "share the same")

	// p2p with no topology pairs.
	mustErr("p2p without topology pairs", twoTierDevices(), &Config{
		P2p: &P2pConfig{},
	}, "no link pairs")
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
