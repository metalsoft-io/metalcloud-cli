package fabric_switch_config

// NVIDIA tags that drive the ASN / loopback / topology sort and grouping.
const (
	tagPod        = "nvidia/pod-id"
	tagSu         = "nvidia/scalability-unit-id"
	tagRail       = "nvidia/rail-group-id"
	tagSpineIndex = "nvidia/spine-index"
	tagSspGroup   = "nvidia/ssp-group-id"
)

// ASN starting points per role (private 32-bit ASN space), overridable in the
// config's asn section.
const (
	defaultLeafStart       int64 = 4_200_000_000
	defaultSpineStart      int64 = 4_201_000_000
	defaultSuperSpineStart int64 = 4_202_000_000
	maxASN                 int64 = 4_294_967_294
)

const defaultLoopbackSubnet = "10.253.128.0/18"

// Reference port layout. All fabric ports are 2x-breakout sub-ports of a 64-OSFP
// switch: logical split l in 1..128 maps to swp{(l+1)//2}s{1 - l%2}.
const (
	splitsPerSwitch           = 128
	leafUplinkLogicalStart    = 65 // swp33s0; leaf splits 1..64 face the hosts
	spineDownlinkLogicalStart = 1  // swp1s0
	spineUplinkLogicalStart   = 33 // swp17s0 (3-tier; downlinks get logicals 1..32)
	sspDownlinkLogicalStart   = 1  // swp1s0
	// Each spine owns a fixed 64-address run of the spine<->superspine pool.
	spineSspRunAddresses = 64
)

// /31 pool defaults per link layer (reference IPAM; overridable via p2p.pools).
var defaultPools = map[string]string{
	"leafSpine":       "10.254.0.0/16",
	"spineSuperSpine": "100.64.0.0/10",
	"leafHost":        "172.16.0.0/12",
}

// The /31 formulas address into the pool by octet position, so each pool must be
// aligned to at least these prefix lengths.
var poolMaxPrefixLen = map[string]int{
	"leafSpine":       24,
	"spineSuperSpine": 26,
	"leafHost":        16,
}

// nvidia/link-layer tag value per topology layer.
var layerTagValue = map[string]string{
	"leafSpine":       "leaf-spine",
	"spineSuperSpine": "spine-superspine",
	"leafHost":        "leaf-server",
}

// Leaf -> host (HGX node) downlink defaults.
var defaultHostNicNames = []string{
	"enp26s0f0np0", "enp60s0f0np0", "enp77s0f0np0", "enp94s0f0np0",
	"enp156s0f0np0", "enp188s0f0np0", "enp204s0f0np0", "enp220s0f0np0",
}

const (
	defaultHostPortPattern       = "swp{port}s{sub}"
	defaultHostDescTemplate      = "to_hgx-su{su:02d}-h{node:02d}_{nic}"
	defaultHostDescTemplate3Tier = "to_hgx-pod{pod:02d}-su{su:02d}-h{node:02d}_{nic}"
	defaultHostNodeCount         = 32
)

// PendingDescription is stamped on physical ports for which no description rule
// exists yet (applied by the runner, not the pure compute step).
const PendingDescription = "dummy_pending_implementation"

// Tags / limits used by the runner when materializing subnets.
const (
	FabricTag        = "metalcloud/fabric-id"
	SubnetNameMaxLen = 63
)

// Built-in hostname templates reproducing the reference naming exactly. The pod
// part exists only in 3-tier. Keyed by three-tier? -> position -> template.
var defaultHostnameTemplates = map[bool]map[string]string{
	false: { // 2-tier: flat fabric, spines carry no pod/rail scope
		"leaf":  "leaf-su{tag:nvidia/scalability-unit-id:02d}-r{tag:nvidia/rail-group-id}",
		"spine": "spine-s{tag:nvidia/spine-index:02d}",
	},
	true: { // 3-tier: pod-scoped; su in the name is the SU within its POD
		"leaf": "leaf-pod{tag:nvidia/pod-id:02d}" +
			"-su{tag:nvidia/scalability-unit-id:02d}-r{tag:nvidia/rail-group-id}",
		"spine": "spine-pod{tag:nvidia/pod-id:02d}" +
			"-r{tag:nvidia/rail-group-id}-s{tag:nvidia/spine-index:02d}",
		"super_spine": "ssp-group{tag:nvidia/ssp-group-id:02d}-s{ordinalBy0:nvidia/ssp-group-id:02d}",
	},
}
