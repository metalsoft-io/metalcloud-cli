package fabric_switch_config

import (
	"gopkg.in/yaml.v3"
)

// maxHostNodes is the number of leaf swpNs0/s1 pairs below the uplink block.
const maxHostNodes = (leafUplinkLogicalStart - 1) / 2 // 32

// rawConfig is the YAML shape. Pointer / map fields let us distinguish an absent
// section (nil) from a present-but-empty one (non-nil), which is what enables a
// feature ("asn: {}" => ASNs with defaults). Unknown top-level keys (api,
// freeform, bgp, ...) used by the sibling scripts are ignored.
type rawConfig struct {
	Ordering            *string                `yaml:"ordering"`
	Hostname            map[string]*string     `yaml:"hostname"`
	Asn                 map[string]*int64      `yaml:"asn"`
	Loopback            map[string]interface{} `yaml:"loopback"`
	Topology            *rawTopology           `yaml:"topology"`
	P2p                 *rawP2p                `yaml:"p2p"`
	DescriptionTemplate *string                `yaml:"descriptionTemplate"`
	EnablePhysicalPorts *bool                  `yaml:"enablePhysicalPorts"`
}

type rawTopology struct {
	LeafSpine       *rawLayer    `yaml:"leafSpine"`
	SpineSuperSpine *rawLayer    `yaml:"spineSuperSpine"`
	LeafHost        *rawLeafHost `yaml:"leafHost"`
}

type rawLayer struct {
	LinksPerPair interface{} `yaml:"linksPerPair"`
	// Legacy one-port-per-spine keys, rejected with a migration hint.
	LeafUplinkPorts          interface{} `yaml:"leafUplinkPorts"`
	SpineDownlinkPortPattern interface{} `yaml:"spineDownlinkPortPattern"`
}

type rawLeafHost struct {
	NodeCount           *int      `yaml:"nodeCount"`
	Nodes               *[]int    `yaml:"nodes"`
	PortPattern         string    `yaml:"portPattern"`
	NicNames            *[]string `yaml:"nicNames"`
	DescriptionTemplate *string   `yaml:"descriptionTemplate"`
}

type rawP2p struct {
	Pools    map[string]string `yaml:"pools"`
	Mtu      *int32            `yaml:"mtu"`
	Supernet interface{}       `yaml:"supernet"`
}

// LoadConfig parses and validates a fabric-switch configuration from YAML/JSON
// bytes. It mirrors load_config() in configure_switches.py: it validates the
// feature sections' shapes and value ranges and returns a *ConfigError on any
// violation. Structural rules that depend on the device set (tag presence,
// topology fit) are enforced later by ComputeDesired.
func LoadConfig(data []byte) (*Config, error) {
	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, configErrorf("invalid configuration: %s", err.Error())
	}

	config := &Config{
		DescriptionTemplate: raw.DescriptionTemplate,
		EnablePhysicalPorts: raw.EnablePhysicalPorts,
	}

	config.Ordering = OrderingManagementAddress
	if raw.Ordering != nil {
		config.Ordering = *raw.Ordering
	}
	if !isValidOrdering(config.Ordering) {
		return nil, configErrorf("ordering must be one of %v, got %q", validOrderings, config.Ordering)
	}

	if raw.Hostname != nil {
		config.Hostname = &HostnameConfig{Templates: raw.Hostname}
	}

	if raw.Asn != nil {
		asn, err := buildAsn(raw.Asn)
		if err != nil {
			return nil, err
		}
		config.Asn = asn
	}

	if raw.Loopback != nil {
		loopback, err := buildLoopback(raw.Loopback)
		if err != nil {
			return nil, err
		}
		config.Loopback = loopback
	}

	if raw.Topology != nil {
		topo, err := buildTopology(raw.Topology)
		if err != nil {
			return nil, err
		}
		config.Topology = topo
	}

	if config.DescriptionTemplate != nil && config.Topology == nil {
		return nil, configErrorf("'descriptionTemplate' requires a 'topology' section (it describes the link pairs)")
	}

	if raw.P2p != nil {
		p2p, err := buildP2p(raw.P2p, config.Topology)
		if err != nil {
			return nil, err
		}
		config.P2p = p2p
	}

	// Presence (not truthiness) enables a feature.
	anyFeature := config.Hostname != nil || config.Asn != nil || config.Loopback != nil ||
		config.Topology != nil || config.P2p != nil || config.DescriptionTemplate != nil
	if !anyFeature && !config.enablePhysicalPorts() {
		return nil, configErrorf("config enables no feature; nothing to do")
	}

	return config, nil
}

func buildAsn(raw map[string]*int64) (*AsnConfig, error) {
	allowed := map[string]bool{"leafStart": true, "spineStart": true, "superSpineStart": true}
	asn := &AsnConfig{}
	for key, value := range raw {
		if !allowed[key] {
			return nil, configErrorf("unknown asn key %q; allowed: [leafStart spineStart superSpineStart]", key)
		}
		if value == nil || *value < 1 || *value > maxASN {
			return nil, configErrorf("asn.%s must be an integer in [1, %d]", key, maxASN)
		}
		switch key {
		case "leafStart":
			asn.LeafStart = value
		case "spineStart":
			asn.SpineStart = value
		case "superSpineStart":
			asn.SuperSpineStart = value
		}
	}
	return asn, nil
}

func buildLoopback(raw map[string]interface{}) (*LoopbackConfig, error) {
	for key := range raw {
		if key != "subnet" {
			return nil, configErrorf("unknown loopback key %q; allowed: [subnet]", key)
		}
	}
	loopback := &LoopbackConfig{}
	if v, ok := raw["subnet"]; ok && v != nil {
		s, ok := v.(string)
		if !ok {
			return nil, configErrorf("loopback.subnet must be a string")
		}
		loopback.Subnet = s
	}
	subnet := loopback.Subnet
	if subnet == "" {
		subnet = defaultLoopbackSubnet
	}
	network, err := parseIpv4Network(subnet)
	if err != nil {
		return nil, configErrorf("loopback.subnet: invalid network %q: %s", subnet, err.Error())
	}
	if network.prefixLen > 30 {
		return nil, configErrorf("loopback.subnet %s is too small to allocate addresses from", subnet)
	}
	return loopback, nil
}

func buildTopology(raw *rawTopology) (*TopologyConfig, error) {
	topo := &TopologyConfig{}
	var err error
	if raw.LeafSpine != nil {
		if topo.LeafSpine, err = buildLayer(raw.LeafSpine, "leafSpine"); err != nil {
			return nil, err
		}
	}
	if raw.SpineSuperSpine != nil {
		if topo.SpineSuperSpine, err = buildLayer(raw.SpineSuperSpine, "spineSuperSpine"); err != nil {
			return nil, err
		}
	}
	if raw.LeafHost != nil {
		if topo.LeafHost, err = buildLeafHost(raw.LeafHost); err != nil {
			return nil, err
		}
	}
	return topo, nil
}

func buildLayer(raw *rawLayer, name string) (*LayerConfig, error) {
	if raw.LeafUplinkPorts != nil || raw.SpineDownlinkPortPattern != nil {
		return nil, configErrorf(
			"topology.%s uses the old one-port-per-spine model (leafUplinkPorts/spineDownlinkPortPattern); the reference block-port model derives ports from the formulas - the only knob is linksPerPair",
			name)
	}
	layer := &LayerConfig{}
	switch v := raw.LinksPerPair.(type) {
	case nil:
		// "auto"
	case string:
		if v != "auto" {
			return nil, configErrorf("topology.%s.linksPerPair must be 'auto' or a positive integer", name)
		}
	case int:
		if v < 1 {
			return nil, configErrorf("topology.%s.linksPerPair must be 'auto' or a positive integer", name)
		}
		layer.LinksPerPair = ptrIntVal(v)
	case bool:
		return nil, configErrorf("topology.%s.linksPerPair must be 'auto' or a positive integer", name)
	default:
		return nil, configErrorf("topology.%s.linksPerPair must be 'auto' or a positive integer", name)
	}
	return layer, nil
}

func buildLeafHost(raw *rawLeafHost) (*LeafHostConfig, error) {
	lh := &LeafHostConfig{
		PortPattern:         raw.PortPattern,
		DescriptionTemplate: raw.DescriptionTemplate,
	}
	if raw.Nodes != nil {
		if raw.NodeCount != nil {
			return nil, configErrorf("topology.leafHost.nodes and nodeCount are mutually exclusive (nodes lists the exact 0-based node indices)")
		}
		nodes := *raw.Nodes
		if len(nodes) == 0 {
			return nil, configErrorf("topology.leafHost.nodes must be a non-empty list of integers in [0, %d]", maxHostNodes-1)
		}
		seen := map[int]bool{}
		for _, n := range nodes {
			if n < 0 || n >= maxHostNodes {
				return nil, configErrorf("topology.leafHost.nodes must be a non-empty list of integers in [0, %d]", maxHostNodes-1)
			}
			if seen[n] {
				return nil, configErrorf("topology.leafHost.nodes has duplicate entries")
			}
			seen[n] = true
		}
		lh.Nodes = nodes
	} else if raw.NodeCount != nil {
		if *raw.NodeCount < 1 || *raw.NodeCount > maxHostNodes {
			return nil, configErrorf("topology.leafHost.nodeCount must be an integer in [1, %d]", maxHostNodes)
		}
		lh.NodeCount = raw.NodeCount
	}
	if raw.NicNames != nil {
		nics := *raw.NicNames
		if len(nics) == 0 {
			return nil, configErrorf("topology.leafHost.nicNames must be a non-empty list of strings")
		}
		if len(nics)%2 != 0 {
			return nil, configErrorf("topology.leafHost.nicNames must have an even number of entries (s0 uses nic[railGroup], s1 uses nic[railGroup + half])")
		}
		lh.NicNames = nics
	}
	return lh, nil
}

func buildP2p(raw *rawP2p, topo *TopologyConfig) (*P2pConfig, error) {
	if topo == nil {
		return nil, configErrorf("'p2p' requires a 'topology' section defining the link pairs")
	}
	if raw.Supernet != nil {
		return nil, configErrorf("p2p.supernet (sequential carving) is replaced by the deterministic per-layer reference formulas; configure p2p.pools.{leafSpine, spineSuperSpine, leafHost} instead")
	}
	p2p := &P2pConfig{Mtu: raw.Mtu}
	if raw.Pools != nil {
		for layer := range raw.Pools {
			if _, ok := defaultPools[layer]; !ok {
				return nil, configErrorf("unknown p2p.pools key %q; allowed: [leafHost leafSpine spineSuperSpine]", layer)
			}
		}
		p2p.Pools = raw.Pools
	}
	// Validate each pool (configured or default) is an aligned IPv4 network.
	for layer, def := range defaultPools {
		raw := def
		if v, ok := p2p.Pools[layer]; ok && v != "" {
			raw = v
		}
		pool, err := parseIpv4Network(raw)
		if err != nil {
			return nil, configErrorf("p2p.pools.%s: invalid network %q: %s", layer, raw, err.Error())
		}
		if pool.prefixLen > poolMaxPrefixLen[layer] {
			return nil, configErrorf(
				"p2p.pools.%s must be /%d or larger (the %s /31 formula addresses into the pool by octet position), got %q",
				layer, poolMaxPrefixLen[layer], layerTagValue[layer], raw)
		}
	}
	return p2p, nil
}

func ptrIntVal(v int) *int { return &v }
