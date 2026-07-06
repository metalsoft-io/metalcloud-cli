package fabric_switch_config

import "testing"

// fakeClient is an in-memory Client that records writes and applies them so a
// re-run observes the new state.
type fakeClient struct {
	fabric      *FabricInfo
	devices     map[int64]*DeviceRecord
	ports       map[int64][]*PortRecord
	p2pLinks    []*P2pLinkRecord
	subnets     []*SubnetRecord
	siteDevices []*DeviceRecord
	nextId      int64

	devicePatches   int
	portPatches     int
	portIpAdds      int
	linksCreated    []P2pLinkCreate
	subnetsCreated  []SubnetCreate
	strategyCreates int
	siteListCalls   int
}

func (f *fakeClient) newId() int64 { f.nextId++; return f.nextId }

func (f *fakeClient) GetFabric(int64) (*FabricInfo, error) { return f.fabric, nil }

func (f *fakeClient) ListFabricDevices(int64) ([]*DeviceRecord, error) {
	out := make([]*DeviceRecord, 0, len(f.devices))
	for _, d := range f.devices {
		out = append(out, d)
	}
	return out, nil
}

func (f *fakeClient) ListDevicesBySite(int64) ([]*DeviceRecord, error) {
	f.siteListCalls++
	return f.siteDevices, nil
}

func (f *fakeClient) UpdateDevice(deviceId int64, body DeviceUpdate, _ int64) error {
	f.devicePatches++
	d := f.devices[deviceId]
	if body.IdentifierString != nil {
		d.IdentifierString = *body.IdentifierString
	}
	if body.ApplyIdentifierAsHostnameOnNextDeploy != nil {
		d.ApplyIdentifierAsHostnameOnNextDeploy = *body.ApplyIdentifierAsHostnameOnNextDeploy
	}
	if body.Asn != nil {
		d.Asn = *body.Asn
	}
	if body.LoopbackAddress != nil {
		d.LoopbackAddressIpv4 = body.LoopbackAddress
	}
	d.Revision++
	return nil
}

func (f *fakeClient) ListPorts(deviceId int64) ([]*PortRecord, error) { return f.ports[deviceId], nil }

func (f *fakeClient) UpdatePortConfig(deviceId, portId int64, enabled *bool, description *string, _ int64) error {
	f.portPatches++
	for _, p := range f.ports[deviceId] {
		if p.InterfaceId == portId {
			if enabled != nil {
				p.Enabled = enabled
			}
			if description != nil {
				p.Description = description
			}
			p.ConfigRevision++
		}
	}
	return nil
}

func (f *fakeClient) AddPortIpv4(deviceId, portId int64, address string, prefixLength int32, _ int64) error {
	f.portIpAdds++
	for _, p := range f.ports[deviceId] {
		if p.InterfaceId == portId {
			p.Ipv4Addresses = append(p.Ipv4Addresses, IpAddress{Address: address, PrefixLength: prefixLength})
		}
	}
	return nil
}

func (f *fakeClient) ListP2pLinks() ([]*P2pLinkRecord, error) { return f.p2pLinks, nil }

func (f *fakeClient) CreateP2pLink(payload P2pLinkCreate) (*P2pLinkRecord, error) {
	f.linksCreated = append(f.linksCreated, payload)
	link := &P2pLinkRecord{Id: f.newId(), Revision: 1}
	a := payload.InterfaceAId
	link.InterfaceAId = &a
	if payload.InterfaceBId != nil {
		b := *payload.InterfaceBId
		link.InterfaceBId = &b
	}
	link.HasIpv4Strategy = payload.StagedSubnetId != nil
	f.p2pLinks = append(f.p2pLinks, link)
	return link, nil
}

func (f *fakeClient) CreateP2pIpv4Strategy(linkId, _ int64, _ string, _ int64) error {
	f.strategyCreates++
	for _, l := range f.p2pLinks {
		if l.Id == linkId {
			l.HasIpv4Strategy = true
			l.Revision++
		}
	}
	return nil
}

func (f *fakeClient) ListSubnetsByFabricTag(int64) ([]*SubnetRecord, error) { return f.subnets, nil }

func (f *fakeClient) CreateSubnet(payload SubnetCreate) (*SubnetRecord, error) {
	f.subnetsCreated = append(f.subnetsCreated, payload)
	s := &SubnetRecord{Id: f.newId(), NetworkAddress: payload.NetworkAddress, PrefixLength: payload.PrefixLength, Tags: payload.Tags}
	f.subnets = append(f.subnets, s)
	return s, nil
}

// newFakeFromFixture builds a fake whose ports cover every port referenced by
// the computed plan (so links resolve), plus a loopback per device and one
// spare physical port per device to exercise the placeholder rule.
func newFakeFromFixture(t *testing.T, devices []*Device, config *Config) (*fakeClient, *DesiredState) {
	t.Helper()
	groups, err := GroupAndOrder(devices, config.ordering())
	if err != nil {
		t.Fatalf("GroupAndOrder: %v", err)
	}
	state, err := ComputeDesired(config, groups)
	if err != nil {
		t.Fatalf("ComputeDesired: %v", err)
	}

	// Collect port names per device from the plan.
	portNames := map[int64]map[string]bool{}
	add := func(id int64, name string) {
		if portNames[id] == nil {
			portNames[id] = map[string]bool{}
		}
		portNames[id][name] = true
	}
	for _, l := range state.Links {
		add(l.DeviceA.Id, l.PortA)
		add(l.DeviceB.Id, l.PortB)
	}
	for _, h := range state.HostLinks {
		add(h.Leaf.Id, h.LeafPort)
	}
	for k := range state.PortDescriptions {
		add(k.DeviceId, k.PortName)
	}

	f := &fakeClient{
		fabric:  &FabricInfo{Id: 5, Name: "Test Fabric", SiteId: ptrInt64(11)},
		devices: map[int64]*DeviceRecord{},
		ports:   map[int64][]*PortRecord{},
		nextId:  10000,
	}
	for _, d := range devices {
		f.devices[d.Id] = &DeviceRecord{Device: *d, Revision: 1}
		var ports []*PortRecord
		for name := range portNames[d.Id] {
			ports = append(ports, &PortRecord{InterfaceId: f.newId(), InterfaceName: name, Kind: "physical", ConfigRevision: 1})
		}
		// one uncovered spare physical port (placeholder rule) + a loopback
		ports = append(ports, &PortRecord{InterfaceId: f.newId(), InterfaceName: "swp64s1-spare", Kind: "physical", ConfigRevision: 1})
		ports = append(ports, &PortRecord{InterfaceId: f.newId(), InterfaceName: "lo", Kind: "loopback", ConfigRevision: 1})
		f.ports[d.Id] = ports
	}
	return f, state
}

func ptrInt64(v int64) *int64 { return &v }

func TestRunnerWriteRunIdempotencyAndDryRun(t *testing.T) {
	f, state := newFakeFromFixture(t, fixtureDevices(), fixtureConfig())

	res, err := Configure(f, fixtureConfig(), 5, false)
	if err != nil {
		t.Fatalf("Configure: %v", err)
	}
	if res.Failures != 0 {
		t.Fatalf("failures = %d, want 0", res.Failures)
	}

	if f.devicePatches != 8 {
		t.Errorf("device patches = %d, want 8", f.devicePatches)
	}
	fabricLinks, hostLinks := len(state.Links), len(state.HostLinks)
	if got := len(f.linksCreated); got != fabricLinks+hostLinks {
		t.Errorf("links created = %d, want %d", got, fabricLinks+hostLinks)
	}
	if got := len(f.subnetsCreated); got != fabricLinks+hostLinks {
		t.Errorf("subnets created = %d, want %d", got, fabricLinks+hostLinks)
	}
	if res.Counters["/31 strategies added"] != fabricLinks+hostLinks {
		t.Errorf("strategies staged = %d, want %d", res.Counters["/31 strategies added"], fabricLinks+hostLinks)
	}
	if res.Counters["loopback IPs added"] != 8 {
		t.Errorf("loopback IPs added = %d, want 8", res.Counters["loopback IPs added"])
	}
	if f.strategyCreates != 0 {
		t.Errorf("config-endpoint strategy POSTs on fresh run = %d, want 0", f.strategyCreates)
	}

	// Bindings: fabric links a_first, host links b_first.
	aFirst, bFirst := 0, 0
	for _, l := range f.linksCreated {
		switch l.StagedBinding {
		case "a_first":
			aFirst++
		case "b_first":
			bFirst++
		}
		if l.RoutingActivation == "while_transporting_logical_network" && l.InterfaceBId != nil {
			t.Errorf("host link should be half-connected (no interfaceB)")
		}
	}
	if aFirst != fabricLinks || bFirst != hostLinks {
		t.Errorf("bindings: a_first=%d b_first=%d, want %d/%d", aFirst, bFirst, fabricLinks, hostLinks)
	}

	// Subnet tags sample (leaf<->spine 10.254.0.0).
	for _, s := range f.subnetsCreated {
		if s.NetworkAddress == "10.254.0.0" {
			if s.Tags["nvidia/link-layer"] != "leaf-spine" || s.Tags[FabricTag] != "5" ||
				s.Tags["nvidia/endpoint-a"] != "leaf-pod5-su1-r1" || s.Tags["nvidia/port-a"] != "swp33s0" {
				t.Errorf("leaf-spine subnet tags wrong: %v", s.Tags)
			}
			if s.PrefixLength != 31 {
				t.Errorf("subnet prefix = %d, want 31", s.PrefixLength)
			}
		}
	}

	// Placeholder description on the uncovered spare physical port.
	if d := portDescription(f, 4, "swp64s1-spare"); d == nil || *d != PendingDescription {
		t.Errorf("spare port placeholder = %v, want %q", d, PendingDescription)
	}

	// ---- idempotency: second run over the mutated fake -> zero writes ----
	f.devicePatches, f.portPatches, f.portIpAdds = 0, 0, 0
	f.linksCreated, f.subnetsCreated, f.strategyCreates = nil, nil, 0
	res2, err := Configure(f, fixtureConfig(), 5, false)
	if err != nil {
		t.Fatalf("second Configure: %v", err)
	}
	if f.devicePatches != 0 || f.portPatches != 0 || f.portIpAdds != 0 ||
		len(f.linksCreated) != 0 || len(f.subnetsCreated) != 0 || f.strategyCreates != 0 {
		t.Errorf("second run made writes: devices=%d ports=%d ips=%d links=%d subnets=%d strategies=%d",
			f.devicePatches, f.portPatches, f.portIpAdds, len(f.linksCreated), len(f.subnetsCreated), f.strategyCreates)
	}
	if res2.Failures != 0 {
		t.Errorf("second run failures = %d", res2.Failures)
	}

	// ---- dry-run on a fresh fixture -> zero writes ----
	fdry, _ := newFakeFromFixture(t, fixtureDevices(), fixtureConfig())
	if _, err := Configure(fdry, fixtureConfig(), 5, true); err != nil {
		t.Fatalf("dry-run Configure: %v", err)
	}
	if fdry.devicePatches != 0 || fdry.portPatches != 0 || fdry.portIpAdds != 0 ||
		len(fdry.linksCreated) != 0 || len(fdry.subnetsCreated) != 0 || fdry.strategyCreates != 0 {
		t.Errorf("dry-run made writes")
	}
}

func TestRunnerStrategyRepair(t *testing.T) {
	f, _ := newFakeFromFixture(t, fixtureDevices(), fixtureConfig())
	if _, err := Configure(f, fixtureConfig(), 5, false); err != nil {
		t.Fatalf("Configure: %v", err)
	}
	// Drop the strategy off one existing link; a re-run must repair exactly it.
	f.p2pLinks[0].HasIpv4Strategy = false
	f.devicePatches, f.portPatches, f.portIpAdds = 0, 0, 0
	f.linksCreated, f.subnetsCreated, f.strategyCreates = nil, nil, 0
	if _, err := Configure(f, fixtureConfig(), 5, false); err != nil {
		t.Fatalf("repair Configure: %v", err)
	}
	if f.strategyCreates != 1 {
		t.Errorf("repair strategy POSTs = %d, want 1", f.strategyCreates)
	}
	if f.devicePatches != 0 || len(f.linksCreated) != 0 {
		t.Errorf("repair run made unexpected writes")
	}
}

func TestRunnerTagsHydration(t *testing.T) {
	devices := fixtureDevices()
	site := make([]*DeviceRecord, len(devices))
	for i, d := range devices {
		cp := *d
		site[i] = &DeviceRecord{Device: cp}
	}
	// Strip tags from the fabric listing; the site listing keeps them.
	stripped := make([]*Device, len(devices))
	for i, d := range devices {
		cp := *d
		cp.TagsMap = map[string]string{}
		stripped[i] = &cp
	}
	f, _ := newFakeFromFixture(t, devices, fixtureConfig())
	// Replace device records' tags with empty, set site devices with tags.
	for _, d := range f.devices {
		d.TagsMap = map[string]string{}
	}
	f.siteDevices = site

	res, err := Configure(f, fixtureConfig(), 5, true)
	if err != nil {
		t.Fatalf("Configure with hydration: %v", err)
	}
	if f.siteListCalls != 1 {
		t.Errorf("site list calls = %d, want 1", f.siteListCalls)
	}
	if res.Failures != 0 {
		t.Errorf("failures after hydration = %d", res.Failures)
	}
}

func portDescription(f *fakeClient, deviceId int64, name string) *string {
	for _, p := range f.ports[deviceId] {
		if p.InterfaceName == name {
			return p.Description
		}
	}
	return nil
}
