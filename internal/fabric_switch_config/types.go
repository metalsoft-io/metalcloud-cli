// Package fabric_switch_config is a pure (no I/O, no SDK) port of the compute
// engine of the standalone nvidia-ra-scripts/configure_switches.py script.
//
// It turns a declarative fabric configuration plus the current set of fabric
// devices into a fully-computed DesiredState: per-device hostname/ASN/loopback,
// the point-to-point link plan (with deterministic /31 offsets), the leaf->host
// downlink plan, and the per-port descriptions. Applying that state through the
// SDK (idempotently, with dry-run) is the job of the runner layer (Phase 3);
// nothing in this package performs network calls.
package fabric_switch_config

// Device is the minimal, SDK-independent view of a fabric network device the
// compute engine needs. The runner builds these from sdk.NetworkDevice.
type Device struct {
	Id                int64
	Position          string // "leaf" | "spine" | "super_spine" | ...
	ManagementAddress string
	IdentifierString  string
	Driver            string
	TagsMap           map[string]string
}

// Label returns a stable human-readable identifier for the device, used in
// error messages and as a fallback peer name in descriptions/subnet tags.
func (d *Device) Label() string {
	if d.IdentifierString != "" {
		return d.IdentifierString
	}
	if d.ManagementAddress != "" {
		return d.ManagementAddress
	}
	return "id=" + itoa(d.Id)
}

// DeviceDesired holds the computed per-device target fields. A nil pointer means
// "this feature was not configured for this device" (leave it untouched).
type DeviceDesired struct {
	Hostname   *string
	Asn        *int64
	LoopbackIp *string
}

// Subnet is a computed IPv4 subnet (always a /31 for p2p links).
type Subnet struct {
	NetworkAddress string
	PrefixLength   int
}

func (s Subnet) String() string {
	return s.NetworkAddress + "/" + itoa(int64(s.PrefixLength))
}

// LinkPlan is a fully-connected fabric link. DeviceA is the gateway side
// (interfaceA, smaller/even IP, binding a_first): the leaf on leafSpine links,
// the spine on spineSuperSpine links. PoolOffset is the /31's base-address
// offset from the layer's pool (pool-independent); Subnet is set only once a
// p2p pool has been assigned.
type LinkPlan struct {
	Layer      string // "leafSpine" | "spineSuperSpine"
	DeviceA    *Device
	PortA      string
	DeviceB    *Device
	PortB      string
	PoolOffset int
	Subnet     *Subnet
}

// HostLinkPlan is a leaf->host downlink: a half-connected p2p link (the host
// side is not a MetalSoft switch interface). The /31 inverts the gateway rule:
// the HOST gets the even address and the leaf the odd one (binding b_first),
// with the leaf still on interfaceA.
type HostLinkPlan struct {
	Leaf        *Device
	LeafPort    string
	Description string
	HostName    string // hgx-... node the port serves (for subnet tags)
	Nic         string // remote NIC netdev name (for subnet tags)
	PoolOffset  int
	Subnet      *Subnet
}

// PortKey identifies a (device, port-name) pair.
type PortKey struct {
	DeviceId int64
	PortName string
}

// DesiredState is the full computed plan.
type DesiredState struct {
	ByDevice         map[int64]*DeviceDesired
	Links            []*LinkPlan
	HostLinks        []*HostLinkPlan
	PortDescriptions map[PortKey]string
	Warnings         []string
}

func newDesiredState() *DesiredState {
	return &DesiredState{
		ByDevice:         map[int64]*DeviceDesired{},
		PortDescriptions: map[PortKey]string{},
	}
}

func (s *DesiredState) desiredFor(id int64) *DeviceDesired {
	d, ok := s.ByDevice[id]
	if !ok {
		d = &DeviceDesired{}
		s.ByDevice[id] = d
	}
	return d
}
