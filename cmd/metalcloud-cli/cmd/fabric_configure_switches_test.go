package cmd

import (
	"testing"

	fsc "github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
	"github.com/spf13/cobra"
)

// newConfigureSwitchesTestCmd wires a throwaway command with the same flags as
// the real configure-switches command, so buildSwitchConfigFromFlags can be
// exercised without going through the root command's auth lifecycle.
func newConfigureSwitchesTestCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "configure-switches"}
	cs := &configureSwitchesFlags
	*cs = struct {
		ordering            string
		enablePhysicalPorts bool
		descriptionTemplate string

		hostname           bool
		hostnameLeaf       string
		hostnameSpine      string
		hostnameSuperSpine string
		hostnameSkip       []string

		asn                bool
		asnLeafStart       int64
		asnSpineStart      int64
		asnSuperSpineStart int64

		loopback       bool
		loopbackSubnet string

		topoLeafSpine          bool
		topoLeafSpineLPP       string
		topoSpineSuperSpine    bool
		topoSpineSuperSpineLPP string

		topoLeafHost            bool
		topoLeafHostNodeCount   int
		topoLeafHostNodes       []int
		topoLeafHostPortPattern string
		topoLeafHostNicNames    []string
		topoLeafHostDescription string

		p2p                    bool
		p2pPoolLeafSpine       string
		p2pPoolSpineSuperSpine string
		p2pPoolLeafHost        string
		p2pMtu                 int32
	}{}

	f := cmd.Flags()
	f.StringVar(&cs.ordering, "ordering", "managementAddress", "")
	f.BoolVar(&cs.enablePhysicalPorts, "enable-physical-ports", true, "")
	f.StringVar(&cs.descriptionTemplate, "description-template", "", "")
	f.BoolVar(&cs.hostname, "hostname", false, "")
	f.StringVar(&cs.hostnameLeaf, "hostname-leaf", "", "")
	f.StringVar(&cs.hostnameSpine, "hostname-spine", "", "")
	f.StringVar(&cs.hostnameSuperSpine, "hostname-super-spine", "", "")
	f.StringSliceVar(&cs.hostnameSkip, "hostname-skip", nil, "")
	f.BoolVar(&cs.asn, "asn", false, "")
	f.Int64Var(&cs.asnLeafStart, "asn-leaf-start", 0, "")
	f.Int64Var(&cs.asnSpineStart, "asn-spine-start", 0, "")
	f.Int64Var(&cs.asnSuperSpineStart, "asn-super-spine-start", 0, "")
	f.BoolVar(&cs.loopback, "loopback", false, "")
	f.StringVar(&cs.loopbackSubnet, "loopback-subnet", "", "")
	f.BoolVar(&cs.topoLeafSpine, "topology-leaf-spine", false, "")
	f.StringVar(&cs.topoLeafSpineLPP, "topology-leaf-spine-links-per-pair", "", "")
	f.BoolVar(&cs.topoSpineSuperSpine, "topology-spine-super-spine", false, "")
	f.StringVar(&cs.topoSpineSuperSpineLPP, "topology-spine-super-spine-links-per-pair", "", "")
	f.BoolVar(&cs.topoLeafHost, "topology-leaf-host", false, "")
	f.IntVar(&cs.topoLeafHostNodeCount, "topology-leaf-host-node-count", 0, "")
	f.IntSliceVar(&cs.topoLeafHostNodes, "topology-leaf-host-nodes", nil, "")
	f.StringVar(&cs.topoLeafHostPortPattern, "topology-leaf-host-port-pattern", "", "")
	f.StringSliceVar(&cs.topoLeafHostNicNames, "topology-leaf-host-nic-names", nil, "")
	f.StringVar(&cs.topoLeafHostDescription, "topology-leaf-host-description-template", "", "")
	f.BoolVar(&cs.p2p, "p2p", false, "")
	f.StringVar(&cs.p2pPoolLeafSpine, "p2p-pool-leaf-spine", "", "")
	f.StringVar(&cs.p2pPoolSpineSuperSpine, "p2p-pool-spine-super-spine", "", "")
	f.StringVar(&cs.p2pPoolLeafHost, "p2p-pool-leaf-host", "", "")
	f.Int32Var(&cs.p2pMtu, "p2p-mtu", 0, "")
	return cmd
}

func TestBuildSwitchConfigFromFlags(t *testing.T) {
	cmd := newConfigureSwitchesTestCmd()
	set := map[string]string{
		"hostname":            "true",
		"asn":                 "true",
		"asn-leaf-start":      "4200001000",
		"loopback-subnet":     "10.253.128.0/18",
		"topology-leaf-spine": "true",
		"topology-spine-super-spine-links-per-pair": "4",
		"topology-leaf-host-node-count":             "2",
		"p2p":                                       "true",
		"p2p-mtu":                                   "9216",
		"description-template":                      "to_{peerHostname}_{peerPort}",
	}
	for k, v := range set {
		if err := cmd.Flags().Set(k, v); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}

	data, err := buildSwitchConfigFromFlags(cmd)
	if err != nil {
		t.Fatalf("buildSwitchConfigFromFlags: %v", err)
	}
	cfg, err := fsc.LoadConfig(data)
	if err != nil {
		t.Fatalf("LoadConfig of flag-built YAML failed: %v\n---\n%s", err, data)
	}

	if cfg.Hostname == nil {
		t.Error("hostname section should be present")
	}
	if cfg.Asn == nil || cfg.Asn.LeafStart == nil || *cfg.Asn.LeafStart != 4200001000 {
		t.Errorf("asn.leafStart not carried through: %+v", cfg.Asn)
	}
	if cfg.Loopback == nil || cfg.Loopback.Subnet != "10.253.128.0/18" {
		t.Errorf("loopback subnet not carried through: %+v", cfg.Loopback)
	}
	if cfg.Topology == nil || cfg.Topology.LeafSpine == nil {
		t.Error("topology.leafSpine should be present (auto)")
	}
	if cfg.Topology == nil || cfg.Topology.SpineSuperSpine == nil ||
		cfg.Topology.SpineSuperSpine.LinksPerPair == nil || *cfg.Topology.SpineSuperSpine.LinksPerPair != 4 {
		t.Errorf("spineSuperSpine linksPerPair should be 4: %+v", cfg.Topology)
	}
	if cfg.Topology == nil || cfg.Topology.LeafHost == nil ||
		cfg.Topology.LeafHost.NodeCount == nil || *cfg.Topology.LeafHost.NodeCount != 2 {
		t.Errorf("leafHost nodeCount should be 2")
	}
	if cfg.P2p == nil || cfg.P2p.Mtu == nil || *cfg.P2p.Mtu != 9216 {
		t.Errorf("p2p mtu should be 9216")
	}
	if cfg.DescriptionTemplate == nil || *cfg.DescriptionTemplate != "to_{peerHostname}_{peerPort}" {
		t.Errorf("descriptionTemplate not carried through")
	}
}

func TestBuildSwitchConfigFromFlagsEmpty(t *testing.T) {
	cmd := newConfigureSwitchesTestCmd()
	if _, err := buildSwitchConfigFromFlags(cmd); err == nil {
		t.Error("expected error when no configuration flags are set")
	}
}

func TestBuildSwitchConfigHostnameSkipAndAutoLPP(t *testing.T) {
	cmd := newConfigureSwitchesTestCmd()
	_ = cmd.Flags().Set("hostname-spine", "spine-{tag:nvidia/spine-index}")
	_ = cmd.Flags().Set("hostname-skip", "super_spine")
	_ = cmd.Flags().Set("topology-leaf-spine", "true")
	_ = cmd.Flags().Set("topology-leaf-spine-links-per-pair", "auto")

	data, err := buildSwitchConfigFromFlags(cmd)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	cfg, err := fsc.LoadConfig(data)
	if err != nil {
		t.Fatalf("LoadConfig: %v\n%s", err, data)
	}
	if cfg.Hostname == nil {
		t.Fatal("hostname present")
	}
	skip, ok := cfg.Hostname.Templates["super_spine"]
	if !ok || skip != nil {
		t.Errorf("super_spine should be present and null (skip), got ok=%v val=%v", ok, skip)
	}
	if tmpl := cfg.Hostname.Templates["spine"]; tmpl == nil || *tmpl == "" {
		t.Errorf("spine template should be set")
	}
	if cfg.Topology.LeafSpine.LinksPerPair != nil {
		t.Errorf("auto linksPerPair should map to nil, got %v", *cfg.Topology.LeafSpine.LinksPerPair)
	}
}
