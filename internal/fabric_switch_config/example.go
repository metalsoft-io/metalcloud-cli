package fabric_switch_config

// exampleConfigYAML is a ready-to-edit template for `fabric configure-switches`.
// Every feature section is optional - omit one to skip that step. The values
// shown are the reference defaults; the comments explain each knob. It is kept
// in sync with the parser by TestExampleConfigParses.
const exampleConfigYAML = `# Fabric switch configuration for 'metalcloud-cli fabric configure-switches'.
#
# Every feature section is OPTIONAL - omit a section to skip that step entirely.
# Each step is idempotent: current state is read first and only differences are
# written. Run with --dry-run to preview the full computed plan without writing.
#
# Device roles are taken from each device's 'position' (leaf | spine |
# super_spine). A fabric is 3-tier iff it has super_spine devices. Most steps are
# driven by NVIDIA tags on the devices: nvidia/pod-id, nvidia/scalability-unit-id,
# nvidia/rail-group-id, nvidia/spine-index, nvidia/ssp-group-id.

# Stable order defining each device's 1-based ordinal within its position group.
# One of: managementAddress (default) | identifierString | id.
ordering: managementAddress

# Hostname (identifierString) computation. The presence of this section enables
# it; positions not set here use the built-in reference templates for the
# detected tier. Set a position to null to skip it. Placeholders:
#   {tag:<key>[:fmt]}       value from the device's tagsMap ({tag:nvidia/pod-id:02d})
#   {ordinalBy:<key>[:fmt]} 1-based ordinal among same-position devices sharing
#                           the device's value for tag <key>
#   {ordinalBy0:<key>[:fmt]} same, 0-based
#   {ordinal} / {position}  ordinal within the group / the position name
# Devices that get a new hostname also get applyIdentifierAsHostnameOnNextDeploy
# (Cumulus drivers only). Computed hostnames must be unique across the fabric.
hostname: {}
# Explicit overrides per position still win, e.g.:
# hostname:
#   leaf: "leaf-pod{tag:nvidia/pod-id:02d}-su{tag:nvidia/scalability-unit-id:02d}-r{tag:nvidia/rail-group-id}"
#   spine: null   # skip spines

# ASN assignment (presence enables it; all keys optional, defaults shown).
#   leaves:      unique ASNs from leafStart, sorted by (pod, su, rail)
#   spines:      shared per (pod, rail) group from spineStart (2-tier: all share)
#   superspines: all share superSpineStart
asn:
  leafStart: 4200000000
  spineStart: 4201000000
  superSpineStart: 4202000000

# Loopback IPv4 allocation (presence enables it; subnet optional). Addresses are
# allocated incrementally from the start (first device gets .1): leaves, then
# spines, then superspines. Written both as a /32 on the loopback interface and
# into the device-level loopbackAddress field.
loopback:
  subnet: 10.253.128.0/18

# Set enabled: true on the staged config of every physical port (default true).
enablePhysicalPorts: true

# Cabling topology - the shared input for both interface descriptions and
# point-to-point link creation. Ports are derived from the reference block-port
# model; the only per-layer knob is linksPerPair ('auto' or a positive integer).
topology:
  # Leaf <-> spine. 3-tier: a leaf uplinks only to its own (pod, rail) spines;
  # 2-tier: full mesh. auto L = 32 // susPerPod (3-tier) or 128 // leafCount.
  leafSpine:
    linksPerPair: auto
  # Spine <-> superspine (3-TIER ONLY). A spine with spine-index S connects to
  # all superspines of ssp-group S. auto L = 32 // sspsPerGroup (uniform groups).
  spineSuperSpine:
    linksPerPair: auto
  # Leaf -> host (HGX node) downlinks. All keys optional; defaults shown.
  leafHost:
    nodeCount: 32                 # number of host port-pairs per leaf
    # nodes: [0, 8, 16, 24]       # OR the exact 0-based node indices (mutually
    #                             # exclusive with nodeCount)
    portPattern: "swp{port}s{sub}"
    nicNames: [enp26s0f0np0, enp60s0f0np0, enp77s0f0np0, enp94s0f0np0,
               enp156s0f0np0, enp188s0f0np0, enp204s0f0np0, enp220s0f0np0]
    # descriptionTemplate: "to_hgx-su{su:02d}-h{node:02d}_{nic}"   # tier-aware default

# Interface description set on both ends of every topology pair (requires
# topology; works with or without p2p). Placeholders: {peerHostname}, {peerPort}.
# Physical ports covered by no rule get the visible placeholder
# "dummy_pending_implementation".
descriptionTemplate: "to_{peerHostname}_{peerPort}"

# Point-to-point link creation over the topology pairs (requires topology). Each
# link gets a deterministic /31 from its layer's pool, registered as a tagged
# IPAM subnet and attached as a manual allocation strategy.
p2p:
  pools:                          # all optional; reference defaults shown
    leafSpine: 10.254.0.0/16
    spineSuperSpine: 100.64.0.0/10
    leafHost: 172.16.0.0/12
  mtu: 9216                       # optional; applied to created links
`

// ExampleConfigYAML returns a commented, ready-to-edit configuration template
// for `fabric configure-switches`.
func ExampleConfigYAML() string {
	return exampleConfigYAML
}
