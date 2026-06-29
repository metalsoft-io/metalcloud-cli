package fabric_switch_config

import "testing"

func TestLoadConfigYAML(t *testing.T) {
	yaml := `
ordering: managementAddress
hostname: {}
asn: {}
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
descriptionTemplate: "to-{peerHostname}:{peerPort}"
`
	cfg, err := LoadConfig([]byte(yaml))
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Hostname == nil || cfg.Asn == nil || cfg.Loopback == nil ||
		cfg.Topology == nil || cfg.P2p == nil || cfg.DescriptionTemplate == nil {
		t.Fatalf("expected all feature sections present: %+v", cfg)
	}
	if cfg.Topology.LeafSpine == nil || cfg.Topology.LeafSpine.LinksPerPair != nil {
		t.Errorf("leafSpine should be present with auto (nil) linksPerPair")
	}
	if cfg.Topology.SpineSuperSpine == nil || cfg.Topology.SpineSuperSpine.LinksPerPair == nil || *cfg.Topology.SpineSuperSpine.LinksPerPair != 4 {
		t.Errorf("spineSuperSpine linksPerPair should be 4")
	}
	if cfg.P2p.Mtu == nil || *cfg.P2p.Mtu != 9216 {
		t.Errorf("p2p mtu should be 9216")
	}

	// The parsed config drives the engine to the same plan as the programmatic fixture.
	groups, err := GroupAndOrder(fixtureDevices(), cfg.ordering())
	if err != nil {
		t.Fatalf("GroupAndOrder: %v", err)
	}
	state, err := ComputeDesired(cfg, groups)
	if err != nil {
		t.Fatalf("ComputeDesired from YAML: %v", err)
	}
	ls := 0
	for _, l := range state.Links {
		if l.Layer == "leafSpine" {
			ls++
		}
	}
	if ls != 60 {
		t.Errorf("leafSpine links from YAML config = %d, want 60", ls)
	}
}

func TestLoadConfigValidationErrors(t *testing.T) {
	cases := []struct {
		name string
		yaml string
		want string
	}{
		{"bad ordering", "ordering: nope\nasn: {}\n", "ordering must be one of"},
		{"unknown asn key", "asn:\n  bogusStart: 1\n", "unknown asn key"},
		{"asn out of range", "asn:\n  leafStart: 0\n", "must be an integer"},
		{"loopback too small", "loopback:\n  subnet: 10.0.0.0/31\n", "too small"},
		{"descriptionTemplate without topology", "descriptionTemplate: \"x\"\n", "requires a 'topology'"},
		{"p2p without topology", "p2p:\n  mtu: 9000\n", "requires a 'topology'"},
		{"nodes and nodeCount", "topology:\n  leafHost:\n    nodeCount: 2\n    nodes: [0,1]\n", "mutually exclusive"},
		{"nothing to do", "enablePhysicalPorts: false\n", "nothing to do"},
	}
	for _, c := range cases {
		_, err := LoadConfig([]byte(c.yaml))
		if err == nil {
			t.Errorf("%s: expected error", c.name)
			continue
		}
		if !contains(err.Error(), c.want) {
			t.Errorf("%s: error %q does not contain %q", c.name, err.Error(), c.want)
		}
	}
}
