package fabric_switch_config

import (
	"fmt"
	"sort"
)

// ComputeDesired is the pure computation of the full desired state. It returns a
// *ConfigError on any rule violation. groups must come from GroupAndOrder.
func ComputeDesired(config *Config, groups map[string][]*Device) (*DesiredState, error) {
	state := newDesiredState()

	// A fabric is 3-tier iff it has super_spine devices.
	threeTier := len(groups["super_spine"]) > 0

	if err := computeHostnames(config, groups, threeTier, state); err != nil {
		return nil, err
	}
	if err := computeAsns(config, groups, threeTier, state); err != nil {
		return nil, err
	}
	if err := computeLoopbacks(config, groups, threeTier, state); err != nil {
		return nil, err
	}
	if err := computeTopology(config, groups, threeTier, state); err != nil {
		return nil, err
	}
	if err := assignP2pSubnets(config, state); err != nil {
		return nil, err
	}
	return state, nil
}

// ---- Hostnames --------------------------------------------------------------

func computeHostnames(config *Config, groups map[string][]*Device, threeTier bool, state *DesiredState) error {
	if config.Hostname == nil {
		return nil
	}
	defaults := defaultHostnameTemplates[threeTier]

	// Union of group positions and explicitly-configured positions.
	positions := map[string]bool{}
	for p := range groups {
		positions[p] = true
	}
	for p := range config.Hostname.Templates {
		positions[p] = true
	}

	hostnameTemplates := map[string]string{}
	for p := range positions {
		var template string
		if t, ok := config.Hostname.Templates[p]; ok {
			if t != nil {
				template = *t
			}
		} else {
			template = defaults[p]
		}
		if template != "" {
			hostnameTemplates[p] = template
		}
	}

	hostnameOwners := map[string]*Device{}
	// Deterministic position order for stable error reporting.
	for _, position := range sortedKeys(hostnameTemplates) {
		template := hostnameTemplates[position]
		devs := groups[position]

		subgroupCache := map[string]map[int64]int{}
		subgroupOrdinals := func(key string) (map[int64]int, error) {
			if cached, ok := subgroupCache[key]; ok {
				return cached, nil
			}
			counts := map[string]int{}
			ordinals := map[int64]int{}
			for _, member := range devs {
				value, ok := member.TagsMap[key]
				if !ok {
					return nil, configErrorf("tag '%s' not present in tagsMap of %s", key, member.Label())
				}
				counts[value]++
				ordinals[member.Id] = counts[value]
			}
			subgroupCache[key] = ordinals
			return ordinals, nil
		}

		for idx, dev := range devs {
			dev := dev
			ordinalBy := func(key string) (int, error) {
				ords, err := subgroupOrdinals(key)
				if err != nil {
					return 0, err
				}
				return ords[dev.Id], nil
			}
			values := map[string]any{"ordinal": idx + 1, "position": position}
			hostname, err := expandTemplate(template, dev.TagsMap, ordinalBy, values)
			if err != nil {
				return configErrorf("hostname for %s: %s", dev.Label(), err.Error())
			}
			if owner, exists := hostnameOwners[hostname]; exists {
				return configErrorf(
					"duplicate computed hostname %q for %s and %s (check the devices' tags)",
					hostname, owner.Label(), dev.Label())
			}
			hostnameOwners[hostname] = dev
			h := hostname
			state.desiredFor(dev.Id).Hostname = &h
		}
	}
	return nil
}

// ---- ASNs -------------------------------------------------------------------

func computeAsns(config *Config, groups map[string][]*Device, threeTier bool, state *DesiredState) error {
	if config.Asn == nil {
		return nil
	}
	leafStart, spineStart, superSpineStart := defaultLeafStart, defaultSpineStart, defaultSuperSpineStart
	if config.Asn.LeafStart != nil {
		leafStart = *config.Asn.LeafStart
	}
	if config.Asn.SpineStart != nil {
		spineStart = *config.Asn.SpineStart
	}
	if config.Asn.SuperSpineStart != nil {
		superSpineStart = *config.Asn.SuperSpineStart
	}

	assign := func(dev *Device, asn int64) error {
		if asn > maxASN {
			return configErrorf("computed ASN %d for %s exceeds %d", asn, dev.Label(), maxASN)
		}
		a := asn
		state.desiredFor(dev.Id).Asn = &a
		return nil
	}

	// Leaves: unique ASNs by (pod,su,rail) sort (3-tier) or (su,rail) (2-tier).
	leafSortTags := []string{tagSu, tagRail}
	if threeTier {
		leafSortTags = []string{tagPod, tagSu, tagRail}
	}
	leaves, err := sortByNumericTags(groups["leaf"], leafSortTags)
	if err != nil {
		return err
	}
	for offset, dev := range leaves {
		if err := assign(dev, leafStart+int64(offset)); err != nil {
			return err
		}
	}

	// Spines.
	if threeTier {
		spines, err := sortByNumericTags(groups["spine"], []string{tagPod, tagRail})
		if err != nil {
			return err
		}
		groupIndex := -1
		var previousKey []int
		for _, dev := range spines {
			key, _ := numericTags(dev, tagPod, tagRail)
			if previousKey == nil || !equalIntSlice(key, previousKey) {
				groupIndex++
				previousKey = key
			}
			if err := assign(dev, spineStart+int64(groupIndex)); err != nil {
				return err
			}
		}
	} else {
		for _, dev := range groups["spine"] {
			if err := assign(dev, spineStart); err != nil {
				return err
			}
		}
	}

	// Superspines: one shared ASN.
	for _, dev := range groups["super_spine"] {
		if err := assign(dev, superSpineStart); err != nil {
			return err
		}
	}
	return nil
}

// ---- Loopbacks --------------------------------------------------------------

func computeLoopbacks(config *Config, groups map[string][]*Device, threeTier bool, state *DesiredState) error {
	if config.Loopback == nil {
		return nil
	}
	subnetStr := config.Loopback.Subnet
	if subnetStr == "" {
		subnetStr = defaultLoopbackSubnet
	}
	network, err := parseIpv4Network(subnetStr)
	if err != nil {
		return configErrorf("loopback.subnet: invalid network %q: %s", subnetStr, err.Error())
	}

	leafKeys := []string{tagSu, tagRail}
	spineKeys := []string{tagSpineIndex}
	if threeTier {
		leafKeys = []string{tagPod, tagSu, tagRail}
		spineKeys = []string{tagPod, tagRail, tagSpineIndex}
	}

	var ordered []*Device
	leaves, err := sortByNumericTags(groups["leaf"], leafKeys)
	if err != nil {
		return err
	}
	ordered = append(ordered, leaves...)
	spines, err := sortByNumericTags(groups["spine"], spineKeys)
	if err != nil {
		return err
	}
	ordered = append(ordered, spines...)
	// Stable sort by ssp-group: within a group the fabric ordering holds.
	ssps, err := sortByNumericTags(groups["super_spine"], []string{tagSspGroup})
	if err != nil {
		return err
	}
	ordered = append(ordered, ssps...)

	broadcast := network.lastAddress()
	for offset, dev := range ordered {
		address := network.base + 1 + uint32(offset)
		if address >= broadcast {
			return configErrorf(
				"loopback.subnet %s exhausted at %s (device %d of %d)",
				subnetStr, dev.Label(), offset+1, len(ordered))
		}
		ip := uintToIpv4(address)
		state.desiredFor(dev.Id).LoopbackIp = &ip
	}
	return nil
}

// ---- Topology + descriptions ------------------------------------------------

func computeTopology(config *Config, groups map[string][]*Device, threeTier bool, state *DesiredState) error {
	topo := config.Topology
	if topo == nil {
		return nil
	}

	// Global spine index = ordinal in the numeric sort by (pod,rail,spine-index)
	// (just (spine-index) in 2-tier). Needed by both fabric layers.
	var spinesSorted []*Device
	spineGlobal := map[int64]int{}
	if topo.LeafSpine != nil || topo.SpineSuperSpine != nil {
		spineSortTags := []string{tagSpineIndex}
		if threeTier {
			spineSortTags = []string{tagPod, tagRail, tagSpineIndex}
		}
		var err error
		spinesSorted, err = sortByNumericTags(groups["spine"], spineSortTags)
		if err != nil {
			return err
		}
		seen := map[string]*Device{}
		for _, spine := range spinesSorted {
			key, _ := numericTags(spine, spineSortTags...)
			ks := fmt.Sprint(key)
			if other, ok := seen[ks]; ok {
				return configErrorf(
					"spines %s and %s share the same (%s) tag values %v",
					other.Label(), spine.Label(), joinStrings(spineSortTags, ", "), key)
			}
			seen[ks] = spine
		}
		for i, spine := range spinesSorted {
			spineGlobal[spine.Id] = i
		}
	}

	if err := computeLeafSpine(topo, groups, threeTier, spinesSorted, spineGlobal, state); err != nil {
		return err
	}
	if err := computeSpineSuperSpine(topo, groups, threeTier, spinesSorted, spineGlobal, state); err != nil {
		return err
	}
	if err := computeDescriptions(config, state); err != nil {
		return err
	}
	if err := computeLeafHost(topo, groups, threeTier, state); err != nil {
		return err
	}
	return nil
}

func computeLeafSpine(topo *TopologyConfig, groups map[string][]*Device, threeTier bool, spinesSorted []*Device, spineGlobal map[int64]int, state *DesiredState) error {
	leafSpine := topo.LeafSpine
	if leafSpine == nil {
		return nil
	}
	leaves := groups["leaf"]
	if len(leaves) == 0 {
		return configErrorf("topology.leafSpine configured but the fabric has no leaf devices")
	}
	if len(spinesSorted) == 0 {
		return configErrorf("topology.leafSpine configured but the fabric has no spine devices")
	}

	leafUplinkBudget := splitsPerSwitch - (leafUplinkLogicalStart - 1) // 64
	spineDownlinkBudget := splitsPerSwitch                             // 2-tier
	if threeTier {
		spineDownlinkBudget = spineUplinkLogicalStart - spineDownlinkLogicalStart // 32
	}

	type block struct {
		label  string
		leaves []*Device
		spines []*Device
	}
	var blocks []block
	var autoBlockSize int

	if threeTier {
		// Group leaves by (pod, rail), ordered within block by su.
		leavesSorted, err := sortByNumericTags(leaves, []string{tagPod, tagRail, tagSu})
		if err != nil {
			return err
		}
		leavesByRail, railOrder, err := groupByTagPair(leavesSorted, tagPod, tagRail)
		if err != nil {
			return err
		}
		// Reference invariants.
		suSetsByPod := map[int][]string{} // pod -> distinct su-set signatures
		suSetSeenByPod := map[int]map[string]bool{}
		podSuCount := map[int]int{}
		for _, key := range railOrder {
			railLeaves := leavesByRail[key.s]
			var sus []int
			seenSu := map[int]bool{}
			for _, leaf := range railLeaves {
				su, _ := numericTag(leaf, tagSu)
				if seenSu[su] {
					return configErrorf(
						"duplicate nvidia/scalability-unit-id among the leaves of pod %d rail-group %d: %v",
						key.a, key.b, intsOf(railLeaves, tagSu))
				}
				seenSu[su] = true
				sus = append(sus, su)
			}
			sig := fmt.Sprint(sus)
			if suSetSeenByPod[key.a] == nil {
				suSetSeenByPod[key.a] = map[string]bool{}
			}
			if !suSetSeenByPod[key.a][sig] {
				suSetSeenByPod[key.a][sig] = true
				suSetsByPod[key.a] = append(suSetsByPod[key.a], sig)
			}
			podSuCount[key.a] = len(sus)
		}
		for pod, sigs := range suSetsByPod {
			if len(sigs) > 1 {
				return configErrorf(
					"the rail-groups of pod %d serve different SU sets; the reference fabric has one leaf per (pod, su, rail-group)",
					pod)
			}
		}
		var counts []int
		for _, c := range podSuCount {
			counts = append(counts, c)
		}
		if !allEqual(counts) {
			return configErrorf(
				"pods have different SU counts %v; the reference addressing (L = 32 // susPerPod) needs a uniform SU count",
				podSuCount)
		}
		autoBlockSize = counts[0] // susPerPod

		spinesByRail, spineRailOrder, err := groupByTagPair(spinesSorted, tagPod, tagRail)
		if err != nil {
			return err
		}
		for _, key := range railOrder {
			railSpines := spinesByRail[key.s]
			if len(railSpines) == 0 {
				return configErrorf(
					"no spines with pod %d / rail-group %d for leaf %s",
					key.a, key.b, leavesByRail[key.s][0].Label())
			}
			blocks = append(blocks, block{
				label:  fmt.Sprintf("pod %d rail-group %d", key.a, key.b),
				leaves: leavesByRail[key.s],
				spines: railSpines,
			})
		}
		// Spines whose (pod,rail) has no leaves -> warning, no links.
		for _, key := range spineRailOrder {
			if _, ok := leavesByRail[key.s]; !ok {
				state.warn(fmt.Sprintf("spines of pod %d rail-group %d have no leaves; no links planned", key.a, key.b))
			}
		}
	} else {
		leavesSorted, err := sortByNumericTags(leaves, []string{tagSu, tagRail})
		if err != nil {
			return err
		}
		seen := map[string]*Device{}
		for _, leaf := range leavesSorted {
			key, _ := numericTags(leaf, tagSu, tagRail)
			ks := fmt.Sprint(key)
			if other, ok := seen[ks]; ok {
				return configErrorf(
					"leaves %s and %s share the same (scalability-unit, rail-group) tag values %v",
					other.Label(), leaf.Label(), key)
			}
			seen[ks] = leaf
		}
		blocks = append(blocks, block{label: "the full-mesh fabric", leaves: leavesSorted, spines: spinesSorted})
		autoBlockSize = len(leavesSorted)
	}

	L := 0
	if leafSpine.LinksPerPair == nil { // "auto"
		if autoBlockSize == 0 {
			return configErrorf("topology.leafSpine: no leaves to size links per pair")
		}
		L = spineDownlinkBudget / autoBlockSize
		if L < 1 {
			return configErrorf(
				"%d leaves per spine exceed the spine's %d downlink splits (L would be 0)",
				autoBlockSize, spineDownlinkBudget)
		}
	} else {
		L = *leafSpine.LinksPerPair
	}

	for _, b := range blocks {
		if len(b.spines)*L > leafUplinkBudget {
			return configErrorf(
				"%d spine(s) x %d link(s) for %s exceed the leaf's %d uplink splits",
				len(b.spines), L, b.label, leafUplinkBudget)
		}
		if len(b.leaves)*L > spineDownlinkBudget {
			return configErrorf(
				"%d leaves x %d link(s) for %s exceed the spine's %d downlink splits",
				len(b.leaves), L, b.label, spineDownlinkBudget)
		}
		for leafBlock, leaf := range b.leaves {
			for spineBlock, spine := range b.spines {
				base := spineGlobal[spine.Id]*256 + leafBlock*2*L
				for u := 0; u < L; u++ {
					state.Links = append(state.Links, &LinkPlan{
						Layer:      "leafSpine",
						DeviceA:    leaf,
						PortA:      logicalToSwp(leafUplinkLogicalStart + spineBlock*L + u),
						DeviceB:    spine,
						PortB:      logicalToSwp(spineDownlinkLogicalStart + leafBlock*L + u),
						PoolOffset: base + 2*u,
					})
				}
			}
		}
	}
	return nil
}

func computeSpineSuperSpine(topo *TopologyConfig, groups map[string][]*Device, threeTier bool, spinesSorted []*Device, spineGlobal map[int64]int, state *DesiredState) error {
	spineSsp := topo.SpineSuperSpine
	if spineSsp == nil {
		return nil
	}
	if !threeTier {
		return configErrorf(
			"topology.spineSuperSpine configured but the fabric has no super_spine devices (2-tier fabrics have no superspine layer)")
	}
	if len(spinesSorted) == 0 {
		return configErrorf("topology.spineSuperSpine configured but the fabric has no spine devices")
	}

	// Group superspines by ssp-group, in fabric (group) order.
	sspsByGroup := map[int][]*Device{}
	var groupOrder []int
	for _, ssp := range groups["super_spine"] {
		g, err := numericTag(ssp, tagSspGroup)
		if err != nil {
			return err
		}
		if _, ok := sspsByGroup[g]; !ok {
			groupOrder = append(groupOrder, g)
		}
		sspsByGroup[g] = append(sspsByGroup[g], ssp)
	}

	spineIndexValues := map[int]bool{}
	for _, s := range spinesSorted {
		idx, _ := numericTag(s, tagSpineIndex)
		spineIndexValues[idx] = true
	}
	if !equalIntSets(spineIndexValues, keySet(sspsByGroup)) {
		return configErrorf(
			"nvidia/spine-index values %v and nvidia/ssp-group-id values %v must match exactly (a spine with spine-index S connects to ssp-group S)",
			sortedSet(spineIndexValues), sortedKeysInt(sspsByGroup))
	}

	spineUplinkBudget := spineSspRunAddresses / 2 // 32
	Lss := 0
	if spineSsp.LinksPerPair == nil { // "auto"
		sizes := map[int]bool{}
		for _, g := range sspsByGroup {
			sizes[len(g)] = true
		}
		if len(sizes) != 1 {
			return configErrorf(
				"ssp groups have different sizes; linksPerPair 'auto' (= 32 // sspsPerGroup) needs uniform groups - set topology.spineSuperSpine.linksPerPair explicitly")
		}
		var sspsPerGroup int
		for s := range sizes {
			sspsPerGroup = s
		}
		Lss = spineUplinkBudget / sspsPerGroup
		if Lss < 1 {
			return configErrorf(
				"%d superspines per group exceed the spine's %d uplink splits (L would be 0)",
				sspsPerGroup, spineUplinkBudget)
		}
	} else {
		Lss = *spineSsp.LinksPerPair
	}

	for _, spine := range spinesSorted {
		idx, _ := numericTag(spine, tagSpineIndex)
		groupSsps := sspsByGroup[idx]
		if len(groupSsps)*Lss > spineUplinkBudget {
			return configErrorf(
				"%d superspine(s) x %d link(s) exceed the %d uplink splits of %s",
				len(groupSsps), Lss, spineUplinkBudget, spine.Label())
		}
		sg := spineGlobal[spine.Id]
		if sspDownlinkLogicalStart-1+(sg+1)*Lss > splitsPerSwitch {
			return configErrorf(
				"spine %s (global index %d) x %d link(s) exceed the superspine's %d splits",
				spine.Label(), sg, Lss, splitsPerSwitch)
		}
		for inGroup, ssp := range groupSsps {
			for u := 0; u < Lss; u++ {
				state.Links = append(state.Links, &LinkPlan{
					Layer:      "spineSuperSpine",
					DeviceA:    spine,
					PortA:      logicalToSwp(spineUplinkLogicalStart + inGroup*Lss + u),
					DeviceB:    ssp,
					PortB:      logicalToSwp(sspDownlinkLogicalStart + sg*Lss + u),
					PoolOffset: sg*spineSspRunAddresses + inGroup*2*Lss + 2*u,
				})
			}
		}
	}
	return nil
}

func computeDescriptions(config *Config, state *DesiredState) error {
	if config.DescriptionTemplate == nil || *config.DescriptionTemplate == "" {
		return nil
	}
	template := *config.DescriptionTemplate

	peerHostname := func(dev *Device) string {
		if d, ok := state.ByDevice[dev.Id]; ok && d.Hostname != nil {
			return *d.Hostname
		}
		return dev.Label()
	}

	for _, plan := range state.Links {
		descA, err := expandTemplate(template, nil, nil, map[string]any{
			"peerHostname": peerHostname(plan.DeviceB), "peerPort": plan.PortB,
		})
		if err != nil {
			return err
		}
		state.PortDescriptions[PortKey{plan.DeviceA.Id, plan.PortA}] = descA

		descB, err := expandTemplate(template, nil, nil, map[string]any{
			"peerHostname": peerHostname(plan.DeviceA), "peerPort": plan.PortA,
		})
		if err != nil {
			return err
		}
		state.PortDescriptions[PortKey{plan.DeviceB.Id, plan.PortB}] = descB
	}
	return nil
}

func computeLeafHost(topo *TopologyConfig, groups map[string][]*Device, threeTier bool, state *DesiredState) error {
	leafHost := topo.LeafHost
	if leafHost == nil {
		return nil
	}

	var nodes []int
	if leafHost.Nodes != nil {
		nodes = leafHost.Nodes
	} else {
		nodeCount := defaultHostNodeCount
		if leafHost.NodeCount != nil {
			nodeCount = *leafHost.NodeCount
		}
		for i := 0; i < nodeCount; i++ {
			nodes = append(nodes, i)
		}
	}
	portPattern := leafHost.PortPattern
	if portPattern == "" {
		portPattern = defaultHostPortPattern
	}
	nics := leafHost.NicNames
	if nics == nil {
		nics = defaultHostNicNames
	}
	hostTemplate := defaultHostDescTemplate
	if threeTier {
		hostTemplate = defaultHostDescTemplate3Tier
	}
	if leafHost.DescriptionTemplate != nil {
		hostTemplate = *leafHost.DescriptionTemplate
	}
	half := len(nics) / 2

	// Global SU index = ordinal of the leaf's (pod,su) among all leaves' distinct
	// (pod,su) pairs (just (su) in 2-tier).
	suKeySet := map[string][]int{}
	var suKeys [][]int
	for _, leaf := range groups["leaf"] {
		su, err := numericTag(leaf, tagSu)
		if err != nil {
			return err
		}
		var key []int
		if threeTier {
			pod, err := numericTag(leaf, tagPod)
			if err != nil {
				return err
			}
			key = []int{pod, su}
		} else {
			key = []int{su}
		}
		ks := fmt.Sprint(key)
		if _, ok := suKeySet[ks]; !ok {
			suKeySet[ks] = key
			suKeys = append(suKeys, key)
		}
	}
	sort.SliceStable(suKeys, func(i, j int) bool { return compareIntTuples(suKeys[i], suKeys[j]) })
	if len(suKeys) > 256 {
		return configErrorf(
			"%d scalability units exceed the 256 the host /31 format can encode (third octet = global SU index)",
			len(suKeys))
	}
	suGlobal := map[string]int{}
	for ordinal, key := range suKeys {
		suGlobal[fmt.Sprint(key)] = ordinal
	}

	for _, leaf := range groups["leaf"] {
		su, _ := numericTag(leaf, tagSu)
		railGroup, err := numericTag(leaf, tagRail)
		if err != nil {
			return err
		}
		values := map[string]any{"su": su, "railGroup": railGroup}
		var pod int
		havePod := false
		if threeTier {
			pod, _ = numericTag(leaf, tagPod)
			values["pod"] = pod
			havePod = true
		} else if _, ok := leaf.TagsMap[tagPod]; ok {
			pod, _ = numericTag(leaf, tagPod)
			values["pod"] = pod
			havePod = true
		}
		if railGroup >= half {
			return configErrorf(
				"rail-group index %d of %s is out of range for %d host NIC names (must be < %d)",
				railGroup, leaf.Label(), len(nics), half)
		}
		var suKey []int
		if threeTier {
			suKey = []int{pod, su}
		} else {
			suKey = []int{su}
		}
		suOrdinal := suGlobal[fmt.Sprint(suKey)]

		for _, node := range nodes {
			var hostName string
			if threeTier {
				hostName = fmt.Sprintf("hgx-pod%02d-su%02d-h%02d", pod, su, node)
			} else {
				hostName = fmt.Sprintf("hgx-su%02d-h%02d", su, node)
			}
			for _, sub := range []int{0, 1} {
				portName, err := expandTemplate(portPattern, nil, nil, map[string]any{
					"port": node + 1, "sub": sub, "node": node,
				})
				if err != nil {
					return err
				}
				key := PortKey{leaf.Id, portName}
				if _, exists := state.PortDescriptions[key]; exists {
					return configErrorf(
						"port %s of %s matched by both the leafSpine and leafHost topology; check the port patterns",
						portName, leaf.Label())
				}
				nic := nics[railGroup+sub*half]
				descValues := map[string]any{"node": node, "nic": nic, "su": su, "railGroup": railGroup}
				if havePod {
					descValues["pod"] = pod
				}
				description, err := expandTemplate(hostTemplate, nil, nil, descValues)
				if err != nil {
					return err
				}
				state.PortDescriptions[key] = description
				rail := railGroup + sub*half
				state.HostLinks = append(state.HostLinks, &HostLinkPlan{
					Leaf:        leaf,
					LeafPort:    portName,
					Description: description,
					HostName:    hostName,
					Nic:         nic,
					PoolOffset:  ((2 * rail) << 16) + (suOrdinal << 8) + 2*node,
				})
			}
		}
	}
	return nil
}

// ---- P2P /31 assignment -----------------------------------------------------

func assignP2pSubnets(config *Config, state *DesiredState) error {
	if config.P2p == nil {
		return nil
	}
	if len(state.Links) == 0 && len(state.HostLinks) == 0 {
		return configErrorf("'p2p' configured but 'topology' produced no link pairs")
	}

	pools := map[string]ipv4Network{}
	for layer, def := range defaultPools {
		raw := def
		if config.P2p.Pools != nil {
			if v, ok := config.P2p.Pools[layer]; ok && v != "" {
				raw = v
			}
		}
		pool, err := parseIpv4Network(raw)
		if err != nil {
			return configErrorf("p2p.pools.%s: invalid network %q: %s", layer, raw, err.Error())
		}
		pools[layer] = pool
	}

	assigned := map[string]string{}
	assign := func(layer string, offset int, label string) (*Subnet, error) {
		pool := pools[layer]
		base := pool.base + uint32(offset)
		if !pool.containsSubnet(base, 31) {
			return nil, configErrorf(
				"computed /31 %s/31 for %s falls outside p2p.pools.%s; size the pool to the fabric",
				uintToIpv4(base), label, layer)
		}
		subnet := &Subnet{NetworkAddress: uintToIpv4(base), PrefixLength: 31}
		key := subnet.String()
		if other, ok := assigned[key]; ok {
			return nil, configErrorf(
				"computed /31 %s for %s collides with %s (overlapping p2p.pools?)",
				key, label, other)
		}
		assigned[key] = label
		return subnet, nil
	}

	for _, plan := range state.Links {
		label := fmt.Sprintf("%s:%s<->%s:%s", plan.DeviceA.Label(), plan.PortA, plan.DeviceB.Label(), plan.PortB)
		subnet, err := assign(plan.Layer, plan.PoolOffset, label)
		if err != nil {
			return err
		}
		plan.Subnet = subnet
	}
	for _, plan := range state.HostLinks {
		label := fmt.Sprintf("%s:%s->%s", plan.Leaf.Label(), plan.LeafPort, plan.HostName)
		subnet, err := assign("leafHost", plan.PoolOffset, label)
		if err != nil {
			return err
		}
		plan.Subnet = subnet
	}
	return nil
}
