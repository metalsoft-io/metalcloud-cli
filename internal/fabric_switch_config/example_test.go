package fabric_switch_config

import "testing"

// TestExampleConfigParses guards the shipped example against drifting into an
// invalid shape: it must always load cleanly and enable every feature section.
func TestExampleConfigParses(t *testing.T) {
	cfg, err := LoadConfig([]byte(ExampleConfigYAML()))
	if err != nil {
		t.Fatalf("example config does not parse: %v", err)
	}
	if cfg.Hostname == nil || cfg.Asn == nil || cfg.Loopback == nil ||
		cfg.Topology == nil || cfg.P2p == nil || cfg.DescriptionTemplate == nil {
		t.Errorf("example config should enable all feature sections: %+v", cfg)
	}
	if cfg.Topology.LeafSpine == nil || cfg.Topology.SpineSuperSpine == nil || cfg.Topology.LeafHost == nil {
		t.Errorf("example topology should define all three layers")
	}
}
