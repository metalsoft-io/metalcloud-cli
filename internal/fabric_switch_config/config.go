package fabric_switch_config

// Config is the parsed, validated fabric-switch configuration. Presence (a
// non-nil pointer), not truthiness, enables a feature: an empty Asn{} means
// "ASNs with default starts". The runner builds this from YAML via LoadConfig.
type Config struct {
	Ordering            string
	Hostname            *HostnameConfig
	Asn                 *AsnConfig
	Loopback            *LoopbackConfig
	Topology            *TopologyConfig
	P2p                 *P2pConfig
	DescriptionTemplate *string
	EnablePhysicalPorts *bool // nil => default true
}

// HostnameConfig holds per-position templates. A position present with a nil
// template is an explicit "skip this position"; a position absent falls back to
// the built-in reference template for the detected tier.
type HostnameConfig struct {
	Templates map[string]*string
}

// AsnConfig overrides the per-role ASN starts (nil => default).
type AsnConfig struct {
	LeafStart       *int64
	SpineStart      *int64
	SuperSpineStart *int64
}

// LoopbackConfig sets the pool the loopback /32s are carved from.
type LoopbackConfig struct {
	Subnet string // "" => defaultLoopbackSubnet
}

// TopologyConfig: each non-nil layer enables that layer's pairing.
type TopologyConfig struct {
	LeafSpine       *LayerConfig
	SpineSuperSpine *LayerConfig
	LeafHost        *LeafHostConfig
}

// LayerConfig is the only knob for a fabric layer: links per connected pair.
// LinksPerPair is nil ("auto") or a positive int.
type LayerConfig struct {
	LinksPerPair *int
}

// LeafHostConfig drives the leaf->host downlinks. All fields optional.
type LeafHostConfig struct {
	NodeCount           *int
	Nodes               []int
	PortPattern         string
	NicNames            []string
	DescriptionTemplate *string
}

// P2pConfig enables point-to-point link creation over the topology pairs.
type P2pConfig struct {
	Pools map[string]string // layer -> CIDR ("" entries fall back to defaults)
	Mtu   *int32
}

func (c *Config) ordering() string {
	if c.Ordering == "" {
		return OrderingManagementAddress
	}
	return c.Ordering
}

func (c *Config) enablePhysicalPorts() bool {
	if c.EnablePhysicalPorts == nil {
		return true
	}
	return *c.EnablePhysicalPorts
}
