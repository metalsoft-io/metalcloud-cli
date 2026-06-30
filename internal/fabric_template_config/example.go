package fabric_template_config

// exampleFreeformYAML is a ready-to-edit example for `fabric configure-freeform`.
const exampleFreeformYAML = `# Freeform base-template registration for 'metalcloud-cli fabric configure-freeform'.
#
# This 'freeform' section is combined with the switch-configuration sections
# (topology / p2p / loopback / asn - see 'fabric configure-switches-example'),
# which feed the per-device variables. Run 'configure-switches', a fabric deploy,
# and 'rescan-links' BEFORE this step.
#
# The .j2 template body is uploaded as-is and rendered server-side by the engine.

freeform:
  mode: l3evpn                                  # purel3 | l3evpn (must match bgp.mode)
  templatePath: ./freeform-device-config.j2     # REQUIRED: path to the .j2 template body
  # templateLabel: spectrumx-freeform           # find-or-create idempotency key
  # profilePriority: 50                          # base layer: before BGP 60/61/62
  # applyMode: once                              # once | always
  # hgxPrefix: 172.0.0.0/8                        # default: derived from the leafHost pool
                                                  # (2-tier <oct>.16.0.0/12, 3-tier <oct>.0.0.0/8)

# --- switch-configuration sections this run reads (abbreviated) -----------------
# (include the full sections from 'fabric configure-switches-example')
ordering: managementAddress
loopback:
  subnet: 10.253.128.0/18
topology:
  leafSpine:
    linksPerPair: auto
  leafHost:
    nodeCount: 32
p2p:
  pools:
    leafHost: 172.16.0.0/12
  mtu: 9216
`

// exampleBgpYAML is a ready-to-edit example for `fabric configure-bgp`.
const exampleBgpYAML = `# BGP underlay/overlay/PFC/VRF registration for 'metalcloud-cli fabric configure-bgp'.
#
# Requires the switch-configuration topology.leafSpine + p2p sections (the BGP
# neighbor set IS the link plan) and devices already carrying asn/loopbackAddress
# (run 'fabric configure-switches' first). Combine this 'bgp' section with those.
#
# All four .j2 template paths are required (read up-front); in purel3 only the
# underlay is registered, but the files must still resolve. Bodies are uploaded
# as-is and rendered server-side by the engine.

bgp:
  mode: l3evpn                                       # purel3 | l3evpn (must match freeform.mode)
  templatePath: ./freeform-bgp-underlay.j2           # REQUIRED (underlay)
  # templateLabel: spectrumx-bgp-underlay
  # profilePriority: 60
  # applyMode: once                                  # once | always

  # --- l3evpn overlay RR mesh + PFC + route-domain tenant VRF ---
  overlayTemplatePath: ./freeform-bgp-overlay.j2     # REQUIRED
  # overlayTemplateLabel: spectrumx-bgp-overlay
  # overlayProfilePriority: 61
  pfcTemplatePath: ./freeform-qos-pfc.j2             # REQUIRED
  # pfcTemplateLabel: spectrumx-qos-pfc
  # pfcProfilePriority: 62
  vrfTemplatePath: ./switch-configure-vrf-create.j2  # REQUIRED (action-bound; no profile)
  # vrfTemplateLabel: switch-configure-vrf-create

# --- switch-configuration sections this run reads (abbreviated) -----------------
# (include the full sections from 'fabric configure-switches-example')
ordering: managementAddress
loopback:
  subnet: 10.253.128.0/18
topology:
  leafSpine:
    linksPerPair: auto
  spineSuperSpine:
    linksPerPair: auto
  leafHost:
    nodeCount: 32
p2p:
  pools:
    leafSpine: 10.254.0.0/16
    spineSuperSpine: 100.64.0.0/10
    leafHost: 172.16.0.0/12
  mtu: 9216
`

// ExampleFreeformYAML returns a ready-to-edit configuration template for
// `fabric configure-freeform`.
func ExampleFreeformYAML() string { return exampleFreeformYAML }

// ExampleBgpYAML returns a ready-to-edit configuration template for
// `fabric configure-bgp`.
func ExampleBgpYAML() string { return exampleBgpYAML }
