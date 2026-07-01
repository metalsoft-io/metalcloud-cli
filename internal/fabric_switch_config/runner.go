package fabric_switch_config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
)

// cumulusDrivers are the drivers that support the one-shot
// applyIdentifierAsHostnameOnNextDeploy flag.
var cumulusDrivers = map[string]bool{"cumulus_linux": true, "cumulus42": true}

// RunResult summarizes a Configure run.
type RunResult struct {
	Counters map[string]int
	Failures int
	Warnings []string
}

func (r *RunResult) count(key string)         { r.Counters[key]++ }
func (r *RunResult) countN(key string, n int) { r.Counters[key] += n }

type runner struct {
	client   Client
	config   *Config
	fabricId int64
	dryRun   bool
	state    *DesiredState
	result   *RunResult

	recByID map[int64]*DeviceRecord
	ports   map[int64]map[string]*PortRecord // device id -> port name -> port

	existingLinks map[string]*P2pLinkRecord // unordered iface pair sig -> link
	linksByIface  map[int64]*P2pLinkRecord  // iface id -> link (covers half-connected)
	subnetIds     map[string]int64          // "netaddr/prefix" -> subnet id
}

func (r *runner) fail(format string, args ...any) {
	r.result.Failures++
	logger.Get().Error().Msgf(format, args...)
}

func (r *runner) info(format string, args ...any) {
	logger.Get().Info().Msgf(format, args...)
}

// Configure executes the full configure flow against client. It returns a
// RunResult (counters + failure count); a non-nil error is only returned for
// fatal setup problems (fabric/device fetch, config computation).
func Configure(client Client, config *Config, fabricId int64, dryRun bool) (*RunResult, error) {
	fabric, err := client.GetFabric(fabricId)
	if err != nil {
		return nil, err
	}
	logger.Get().Info().Msgf("Target fabric: %q (id=%d)", fabric.Name, fabricId)

	devices, err := client.ListFabricDevices(fabricId)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("fabric %d has no network devices attached", fabricId)
	}
	logger.Get().Info().Msgf("Fabric has %d device(s)", len(devices))

	if err := hydrateTagsFromSite(client, fabric, devices); err != nil {
		return nil, err
	}

	engineDevices := make([]*Device, len(devices))
	recByID := map[int64]*DeviceRecord{}
	for i, d := range devices {
		engineDevices[i] = &d.Device
		recByID[d.Id] = d
	}

	groups, err := GroupAndOrder(engineDevices, config.ordering())
	if err != nil {
		return nil, err
	}
	state, err := ComputeDesired(config, groups)
	if err != nil {
		return nil, err
	}

	r := &runner{
		client:    client,
		config:    config,
		fabricId:  fabricId,
		dryRun:    dryRun,
		state:     state,
		result:    &RunResult{Counters: map[string]int{}, Warnings: state.Warnings},
		recByID:   recByID,
		ports:     map[int64]map[string]*PortRecord{},
		subnetIds: map[string]int64{},
	}

	logger.Get().Debug().Msgf(
		"Computed plan: %d device target(s), %d fabric link(s), %d host downlink(s), %d port description(s)",
		len(state.ByDevice), len(state.Links), len(state.HostLinks), len(state.PortDescriptions))
	for _, position := range sortedGroupKeys(groups) {
		names := make([]string, 0, len(groups[position]))
		for _, dev := range groups[position] {
			names = append(names, dev.Label())
		}
		logger.Get().Debug().Msgf("position %q: %d device(s): %s", position, len(names), strings.Join(names, ", "))
	}

	targeted := targetedPositions(config)
	targetedNames := make([]string, 0, len(targeted))
	for position := range targeted {
		targetedNames = append(targetedNames, position)
	}
	sort.Strings(targetedNames)
	logger.Get().Debug().Msgf("targeted positions (configure device): %s", strings.Join(targetedNames, ", "))
	for position := range targeted {
		if _, ok := groups[position]; !ok {
			r.result.Warnings = append(r.result.Warnings,
				fmt.Sprintf("config targets position %q but the fabric has no such devices", position))
		}
	}

	processAllPorts := config.enablePhysicalPorts()
	for _, position := range sortedGroupKeys(groups) {
		for _, dev := range groups[position] {
			rec := recByID[dev.Id]
			desired := state.ByDevice[dev.Id]
			if desired == nil {
				desired = &DeviceDesired{}
			}
			if targeted[position] {
				r.configureDevice(rec, desired)
			}
			if targeted[position] || processAllPorts {
				r.configurePorts(rec)
			}
		}
	}

	r.configureLinks()
	return r.result, nil
}

func targetedPositions(config *Config) map[string]bool {
	targeted := map[string]bool{}
	if config.Hostname != nil {
		for position := range config.Hostname.Templates {
			targeted[position] = true
		}
	}
	if config.Asn != nil || config.Loopback != nil {
		targeted["leaf"], targeted["spine"], targeted["super_spine"] = true, true, true
	}
	if config.Topology != nil {
		if config.Topology.LeafSpine != nil {
			targeted["leaf"], targeted["spine"] = true, true
		}
		if config.Topology.SpineSuperSpine != nil {
			targeted["spine"], targeted["super_spine"] = true, true
		}
		if config.Topology.LeafHost != nil {
			targeted["leaf"] = true
		}
	}
	return targeted
}

// hydrateTagsFromSite backfills empty TagsMap fields from the site-scoped device
// listing (the fabric-devices listing can serve {} even when tags exist).
func hydrateTagsFromSite(client Client, fabric *FabricInfo, devices []*DeviceRecord) error {
	allTagged := true
	for _, d := range devices {
		if len(d.TagsMap) == 0 {
			allTagged = false
			break
		}
	}
	if allTagged || fabric.SiteId == nil {
		return nil
	}
	siteDevices, err := client.ListDevicesBySite(*fabric.SiteId)
	if err != nil {
		return err
	}
	tagsById := map[int64]map[string]string{}
	for _, d := range siteDevices {
		if len(d.TagsMap) > 0 {
			tagsById[d.Id] = d.TagsMap
		}
	}
	var hydrated []string
	for _, d := range devices {
		if len(d.TagsMap) == 0 {
			if tags, ok := tagsById[d.Id]; ok {
				d.TagsMap = tags
				hydrated = append(hydrated, d.Label())
			}
		}
	}
	if len(hydrated) > 0 {
		logger.Get().Warn().Msgf(
			"fabric-devices listing served an empty tagsMap for %d device(s); backfilled from the siteId=%d listing",
			len(hydrated), *fabric.SiteId)
	}
	return nil
}

// ---- Devices ----------------------------------------------------------------

func (r *runner) configureDevice(dev *DeviceRecord, desired *DeviceDesired) {
	label := dev.Label()
	body := DeviceUpdate{}

	if desired.Hostname != nil && *desired.Hostname != dev.IdentifierString {
		body.IdentifierString = desired.Hostname
	}
	if desired.Hostname != nil {
		if cumulusDrivers[dev.Driver] {
			if !dev.ApplyIdentifierAsHostnameOnNextDeploy {
				t := true
				body.ApplyIdentifierAsHostnameOnNextDeploy = &t
			}
		} else {
			logger.Get().Warn().Msgf(
				"[%s] driver %q does not support applyIdentifierAsHostnameOnNextDeploy; flag not set",
				label, dev.Driver)
		}
	}
	if desired.Asn != nil && *desired.Asn != dev.Asn {
		body.Asn = desired.Asn
	}
	if desired.LoopbackIp != nil {
		if dev.LoopbackAddressIpv4 == nil || *dev.LoopbackAddressIpv4 != *desired.LoopbackIp {
			body.LoopbackAddress = desired.LoopbackIp
		}
	}

	logger.Get().Debug().Msgf(
		"[%s] device desired: hostname=%s asn=%s loopback=%s (current: hostname=%q asn=%d loopback=%s); patch={%s}",
		label, strOrDash(desired.Hostname), int64OrDash(desired.Asn), strOrDash(desired.LoopbackIp),
		dev.IdentifierString, dev.Asn, strOrDash(dev.LoopbackAddressIpv4), describeDeviceUpdate(body))

	if body.empty() {
		r.result.count("devices unchanged")
		return
	}
	if r.dryRun {
		r.result.count("devices patched")
		r.info("[%s] would PATCH device", label)
		return
	}
	if err := r.client.UpdateDevice(dev.Id, body, dev.Revision); err != nil {
		r.fail("[%s] device PATCH failed: %s", label, err.Error())
		return
	}
	r.result.count("devices patched")
	r.info("[%s] device updated", label)
}

// ---- Ports ------------------------------------------------------------------

func (r *runner) configurePorts(dev *DeviceRecord) {
	label := dev.Label()
	ports, err := r.client.ListPorts(dev.Id)
	if err != nil {
		r.fail("[%s] listing ports failed: %s", label, err.Error())
		return
	}
	byName := map[string]*PortRecord{}
	physical, loopback := 0, 0
	for _, p := range ports {
		if p.InterfaceName != "" {
			byName[p.InterfaceName] = p
		}
		switch p.Kind {
		case "physical":
			physical++
		case "loopback":
			loopback++
		}
	}
	r.ports[dev.Id] = byName
	logger.Get().Debug().Msgf("[%s] %d port(s) discovered: %d physical, %d loopback, %d other",
		label, len(ports), physical, loopback, len(ports)-physical-loopback)
	if loopback == 0 && r.state.ByDevice[dev.Id] != nil && r.state.ByDevice[dev.Id].LoopbackIp != nil {
		logger.Get().Debug().Msgf("[%s] no loopback port among the discovered interfaces; the /32 step will be skipped (have the switch interfaces been discovered yet?)", label)
	}

	enablePhysical := r.config.enablePhysicalPorts()
	haveDescriptionTemplate := r.config.DescriptionTemplate != nil && *r.config.DescriptionTemplate != ""

	for _, port := range ports {
		var enabled *bool
		var description *string

		if enablePhysical && port.Kind == "physical" && (port.Enabled == nil || !*port.Enabled) {
			t := true
			enabled = &t
		}

		desiredDesc, hasDesc := r.state.PortDescriptions[PortKey{dev.Id, port.InterfaceName}]
		if !hasDesc && haveDescriptionTemplate && port.Kind == "physical" {
			desiredDesc, hasDesc = PendingDescription, true
		}
		if hasDesc && (port.Description == nil || *port.Description != desiredDesc) {
			d := desiredDesc
			description = &d
		}

		if enabled == nil && description == nil {
			continue
		}
		logger.Get().Debug().Msgf("[%s:%s] port config patch: enabled=%s description=%s",
			label, port.InterfaceName, boolOrDash(enabled), strOrDash(description))
		if r.dryRun {
			r.result.count("port configs patched")
			continue
		}
		if err := r.client.UpdatePortConfig(dev.Id, port.InterfaceId, enabled, description, port.ConfigRevision); err != nil {
			r.fail("[%s:%s] port config PATCH failed: %s", label, port.InterfaceName, err.Error())
			continue
		}
		r.result.count("port configs patched")
	}

	r.configureLoopbackIp(dev, ports)
}

func (r *runner) configureLoopbackIp(dev *DeviceRecord, ports []*PortRecord) {
	desired := r.state.ByDevice[dev.Id]
	if desired == nil || desired.LoopbackIp == nil {
		return
	}
	label := dev.Label()

	var loopbacks []*PortRecord
	for _, p := range ports {
		if p.Kind == "loopback" {
			loopbacks = append(loopbacks, p)
		}
	}
	if len(loopbacks) == 0 {
		r.fail("[%s] no loopback interface found; cannot set %s/32", label, *desired.LoopbackIp)
		return
	}
	loopback := loopbacks[0]
	for _, p := range loopbacks {
		if p.InterfaceName == "lo" {
			loopback = p
			break
		}
	}
	logger.Get().Debug().Msgf("[%s] loopback interface %q (id=%d) has %d existing address(es); target %s/32",
		label, loopback.InterfaceName, loopback.InterfaceId, len(loopback.Ipv4Addresses), *desired.LoopbackIp)

	for _, addr := range loopback.Ipv4Addresses {
		if addr.Address == *desired.LoopbackIp && addr.PrefixLength == 32 {
			r.result.count("loopback IPs already present")
			return
		}
	}
	if len(loopback.Ipv4Addresses) > 0 {
		r.fail("[%s] loopback %s already has different IP(s); not adding %s/32",
			label, loopback.InterfaceName, *desired.LoopbackIp)
		return
	}
	if r.dryRun {
		r.result.count("loopback IPs added")
		return
	}
	if err := r.client.AddPortIpv4(dev.Id, loopback.InterfaceId, *desired.LoopbackIp, 32, loopback.ConfigRevision); err != nil {
		r.fail("[%s] loopback IP POST failed: %s", label, err.Error())
		return
	}
	r.result.count("loopback IPs added")
}

// ---- P2P links --------------------------------------------------------------

func (r *runner) configureLinks() {
	if r.config.P2p == nil {
		return
	}
	if len(r.state.Links) == 0 && len(r.state.HostLinks) == 0 {
		return
	}

	links, err := r.client.ListP2pLinks()
	if err != nil {
		r.fail("listing point-to-point links failed: %s", err.Error())
		return
	}
	r.existingLinks = map[string]*P2pLinkRecord{}
	r.linksByIface = map[int64]*P2pLinkRecord{}
	for _, link := range links {
		var ids []int64
		if link.InterfaceAId != nil {
			ids = append(ids, *link.InterfaceAId)
			r.linksByIface[*link.InterfaceAId] = link
		}
		if link.InterfaceBId != nil {
			ids = append(ids, *link.InterfaceBId)
			r.linksByIface[*link.InterfaceBId] = link
		}
		if len(ids) == 2 {
			r.existingLinks[ifacePairKey(ids[0], ids[1])] = link
		}
	}
	logger.Get().Debug().Msgf(
		"point-to-point links: %d existing in the system; planning %d fabric link(s) + %d host downlink(s)",
		len(links), len(r.state.Links), len(r.state.HostLinks))

	subnets, err := r.client.ListSubnetsByFabricTag(r.fabricId)
	if err != nil {
		r.fail("listing subnets failed: %s", err.Error())
		return
	}
	for _, s := range subnets {
		r.subnetIds[subnetKey(s.NetworkAddress, int(s.PrefixLength))] = s.Id
	}

	for _, plan := range r.state.Links {
		r.configureLink(plan)
	}
	for _, plan := range r.state.HostLinks {
		r.configureHostLink(plan)
	}
}

func (r *runner) configureLink(plan *LinkPlan) {
	portA := r.resolvePort(plan.DeviceA, plan.PortA)
	portB := r.resolvePort(plan.DeviceB, plan.PortB)
	if portA == nil || portB == nil {
		return
	}
	label := fmt.Sprintf("%s:%s<->%s:%s", plan.DeviceA.Label(), plan.PortA, plan.DeviceB.Label(), plan.PortB)
	tags := r.subnetTags(plan.Layer, r.endpointName(plan.DeviceA), plan.PortA, r.endpointName(plan.DeviceB), plan.PortB)
	name := r.subnetName(plan.DeviceA.Id, plan.PortA, plan.DeviceB.Id, plan.PortB)
	logger.Get().Debug().Msgf("[%s] %s link: ifaceA=%d ifaceB=%d /31=%s binding=%s",
		label, plan.Layer, portA.InterfaceId, portB.InterfaceId, subnetOrDash(plan.Subnet), gatewayBinding)

	if existing, ok := r.existingLinks[ifacePairKey(portA.InterfaceId, portB.InterfaceId)]; ok {
		r.result.count("links existing")
		if plan.Subnet != nil {
			r.ensureLinkStrategy(existing, label, plan.Subnet, tags, gatewayBinding, true, name)
		}
		return
	}

	payload := P2pLinkCreate{
		InterfaceAId:      portA.InterfaceId,
		InterfaceBId:      &portB.InterfaceId,
		Description:       &label,
		Mtu:               r.config.P2p.Mtu,
		RoutingActivation: "default",
	}

	if r.dryRun {
		r.result.count("links created")
		if plan.Subnet != nil {
			r.ensureLinkStrategy(nil, label, plan.Subnet, tags, gatewayBinding, false, name)
		}
		return
	}

	if plan.Subnet != nil {
		if subnetId, ok := r.ensureSubnet(plan.Subnet, tags, name); ok {
			payload.StagedSubnetId = &subnetId
			payload.StagedBinding = gatewayBinding
		}
	}
	if _, err := r.client.CreateP2pLink(payload); err != nil {
		r.fail("[%s] link create failed: %s", label, err.Error())
		return
	}
	r.result.count("links created")
	if payload.StagedSubnetId != nil {
		r.result.count("/31 strategies added")
	}
}

func (r *runner) configureHostLink(plan *HostLinkPlan) {
	port := r.resolvePort(plan.Leaf, plan.LeafPort)
	if port == nil {
		return
	}
	label := fmt.Sprintf("%s:%s (host downlink)", plan.Leaf.Label(), plan.LeafPort)
	tags := r.subnetTags("leafHost", r.endpointName(plan.Leaf), plan.LeafPort, plan.HostName, plan.Nic)
	name := r.hostSubnetName(plan.Leaf.Id, plan.LeafPort, plan.HostName)

	if existing, ok := r.linksByIface[port.InterfaceId]; ok {
		r.result.count("host links existing")
		if plan.Subnet != nil {
			r.ensureLinkStrategy(existing, label, plan.Subnet, tags, hostBinding, true, name)
		}
		return
	}

	payload := P2pLinkCreate{
		InterfaceAId:      port.InterfaceId,
		Description:       &plan.Description,
		Mtu:               r.config.P2p.Mtu,
		RoutingActivation: "while_transporting_logical_network",
	}

	if r.dryRun {
		r.result.count("host links created")
		if plan.Subnet != nil {
			r.ensureLinkStrategy(nil, label, plan.Subnet, tags, hostBinding, false, name)
		}
		return
	}

	if plan.Subnet != nil {
		if subnetId, ok := r.ensureSubnet(plan.Subnet, tags, name); ok {
			payload.StagedSubnetId = &subnetId
			payload.StagedBinding = hostBinding
		}
	}
	if _, err := r.client.CreateP2pLink(payload); err != nil {
		r.fail("[%s] link create failed: %s", label, err.Error())
		return
	}
	r.result.count("host links created")
	if payload.StagedSubnetId != nil {
		r.result.count("/31 strategies added")
	}
}

func (r *runner) resolvePort(dev *Device, portName string) *PortRecord {
	if byName, ok := r.ports[dev.Id]; ok {
		if port, ok := byName[portName]; ok {
			return port
		}
	}
	r.fail("[%s] port %q not found; skipping its link", dev.Label(), portName)
	return nil
}

func (r *runner) ensureSubnet(subnet *Subnet, tags map[string]string, name string) (int64, bool) {
	key := subnetKey(subnet.NetworkAddress, subnet.PrefixLength)
	if id, ok := r.subnetIds[key]; ok {
		logger.Get().Debug().Msgf("subnet %s already exists (id=%d, name=%q)", subnet, id, name)
		return id, true
	}
	logger.Get().Debug().Msgf("subnet %s not found; would create (name=%q, link-layer=%s)", subnet, name, tags["nvidia/link-layer"])
	if r.dryRun {
		r.result.count("subnets created")
		return 0, false
	}
	created, err := r.client.CreateSubnet(SubnetCreate{
		NetworkAddress: subnet.NetworkAddress,
		PrefixLength:   int32(subnet.PrefixLength),
		Name:           name,
		Tags:           tags,
	})
	if err != nil {
		r.fail("creating subnet %s failed: %s", subnet, err.Error())
		return 0, false
	}
	r.subnetIds[key] = created.Id
	r.result.count("subnets created")
	return created.Id, true
}

func (r *runner) ensureLinkStrategy(link *P2pLinkRecord, label string, subnet *Subnet, tags map[string]string, binding string, checkExisting bool, name string) {
	if checkExisting && link != nil && link.HasIpv4Strategy {
		r.result.count("links with existing /31 strategy")
		return
	}
	subnetId, ok := r.ensureSubnet(subnet, tags, name)
	if r.dryRun {
		r.result.count("/31 strategies added")
		return
	}
	if !ok {
		return
	}
	if err := r.client.CreateP2pIpv4Strategy(link.Id, subnetId, binding, link.Revision); err != nil {
		r.fail("[%s] creating /31 strategy failed: %s", label, err.Error())
		return
	}
	r.result.count("/31 strategies added")
}

// ---- subnet tags / names ----------------------------------------------------

const (
	gatewayBinding = "a_first"
	hostBinding    = "b_first"
)

func (r *runner) endpointName(dev *Device) string {
	if d, ok := r.state.ByDevice[dev.Id]; ok && d.Hostname != nil {
		return *d.Hostname
	}
	return dev.Label()
}

func (r *runner) subnetTags(layer, endpointA, portA, endpointB, portB string) map[string]string {
	return map[string]string{
		FabricTag:           fmt.Sprintf("%d", r.fabricId),
		"nvidia/link-layer": layerTagValue[layer],
		"nvidia/endpoint-a": endpointA,
		"nvidia/port-a":     portA,
		"nvidia/endpoint-b": endpointB,
		"nvidia/port-b":     portB,
	}
}

func (r *runner) subnetName(swA int64, portA string, swB int64, portB string) string {
	return truncateName(fmt.Sprintf("fab%d-sw%d-%s-to-sw%d-%s", r.fabricId, swA, portA, swB, portB))
}

func (r *runner) hostSubnetName(leafId int64, leafPort, hostName string) string {
	return truncateName(fmt.Sprintf("fab%d-sw%d-%s-to-%s", r.fabricId, leafId, leafPort, hostName))
}

func truncateName(name string) string {
	if len(name) > SubnetNameMaxLen {
		name = name[:SubnetNameMaxLen]
	}
	for len(name) > 0 && name[len(name)-1] == '-' {
		name = name[:len(name)-1]
	}
	return name
}

// ---- debug formatting helpers -----------------------------------------------

func strOrDash(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

func int64OrDash(n *int64) string {
	if n == nil {
		return "-"
	}
	return strconv.FormatInt(*n, 10)
}

func boolOrDash(b *bool) string {
	if b == nil {
		return "-"
	}
	if *b {
		return "true"
	}
	return "false"
}

func subnetOrDash(s *Subnet) string {
	if s == nil {
		return "-"
	}
	return s.String()
}

// describeDeviceUpdate lists the fields a device patch would change.
func describeDeviceUpdate(body DeviceUpdate) string {
	var parts []string
	if body.IdentifierString != nil {
		parts = append(parts, "identifierString="+*body.IdentifierString)
	}
	if body.ApplyIdentifierAsHostnameOnNextDeploy != nil {
		parts = append(parts, "applyIdentifierAsHostnameOnNextDeploy="+boolOrDash(body.ApplyIdentifierAsHostnameOnNextDeploy))
	}
	if body.Asn != nil {
		parts = append(parts, "asn="+int64OrDash(body.Asn))
	}
	if body.LoopbackAddress != nil {
		parts = append(parts, "loopbackAddress="+*body.LoopbackAddress)
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

// ---- small helpers ----------------------------------------------------------

func ifacePairKey(a, b int64) string {
	if a > b {
		a, b = b, a
	}
	return fmt.Sprintf("%d-%d", a, b)
}

func subnetKey(networkAddress string, prefixLength int) string {
	return fmt.Sprintf("%s/%d", networkAddress, prefixLength)
}

func sortedGroupKeys(groups map[string][]*Device) []string {
	out := make([]string, 0, len(groups))
	for k := range groups {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
