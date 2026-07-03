package fabric_switch_config

// Client is the narrow set of fabric operations the runner needs, abstracted
// away from the SDK so the runner can be unit-tested. sdkClient is the production
// implementation over *sdk.APIClient.
type Client interface {
	GetFabric(fabricId int64) (*FabricInfo, error)
	ListFabricDevices(fabricId int64) ([]*DeviceRecord, error)
	ListDevicesBySite(siteId int64) ([]*DeviceRecord, error)
	UpdateDevice(deviceId int64, body DeviceUpdate, driftStatus string, revision int64) error
	ListPorts(deviceId int64) ([]*PortRecord, error)
	UpdatePortConfig(deviceId, portId int64, enabled *bool, description *string, configRevision int64) error
	AddPortIpv4(deviceId, portId int64, address string, prefixLength int32, configRevision int64) error
	ListP2pLinks() ([]*P2pLinkRecord, error)
	CreateP2pLink(payload P2pLinkCreate) (*P2pLinkRecord, error)
	CreateP2pIpv4Strategy(linkId, subnetId int64, binding string, linkRevision int64) error
	ListSubnetsByFabricTag(fabricId int64) ([]*SubnetRecord, error)
	CreateSubnet(payload SubnetCreate) (*SubnetRecord, error)
}

// FabricInfo is the subset of a fabric the runner reads.
type FabricInfo struct {
	Id     int64
	Name   string
	SiteId *int64
}

// DeviceRecord is a fabric device plus the current-state fields the runner diffs
// against. The embedded Device is what the compute engine consumes.
type DeviceRecord struct {
	Device
	Asn                                   int64
	LoopbackAddressIpv4                   *string
	ApplyIdentifierAsHostnameOnNextDeploy bool
	Revision                              int64
	DriftDetectionSyncStatus              string
}

// DeviceUpdate is a sparse device patch: only the non-nil fields are sent.
type DeviceUpdate struct {
	IdentifierString                      *string
	ApplyIdentifierAsHostnameOnNextDeploy *bool
	Asn                                   *int64
	LoopbackAddress                       *string
}

func (u DeviceUpdate) empty() bool {
	return u.IdentifierString == nil && u.ApplyIdentifierAsHostnameOnNextDeploy == nil &&
		u.Asn == nil && u.LoopbackAddress == nil
}

// PortRecord is a device interface and its staged config.
type PortRecord struct {
	InterfaceId    int64
	InterfaceName  string
	Kind           string
	Enabled        *bool
	Description    *string
	ConfigRevision int64
	Ipv4Addresses  []IpAddress
}

type IpAddress struct {
	Address      string
	PrefixLength int32
}

// P2pLinkRecord is the subset of an existing point-to-point link the runner
// uses for idempotency. InterfaceAId/InterfaceBId are set only for sides of
// type network_equipment_interface.
type P2pLinkRecord struct {
	Id              int64
	Revision        int64
	InterfaceAId    *int64
	InterfaceBId    *int64
	HasIpv4Strategy bool
}

// P2pLinkCreate is a link to create. InterfaceBId nil => half-connected link.
// When StagedSubnetId is non-nil, a manual ipv4 strategy is staged on create.
type P2pLinkCreate struct {
	InterfaceAId      int64
	InterfaceBId      *int64
	Description       *string
	Mtu               *int32
	RoutingActivation string
	StagedSubnetId    *int64
	StagedBinding     string
}

// SubnetRecord is an existing IPAM subnet.
type SubnetRecord struct {
	Id             int64
	NetworkAddress string
	PrefixLength   int32
	Tags           map[string]string
}

// SubnetCreate is a /31 IPAM subnet to create.
type SubnetCreate struct {
	NetworkAddress string
	PrefixLength   int32
	Name           string
	Tags           map[string]string
}
