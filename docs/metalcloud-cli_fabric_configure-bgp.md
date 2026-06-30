## metalcloud-cli fabric configure-bgp

Register the BGP underlay/overlay/PFC templates + profiles (step 8b)

### Synopsis

Register the Spectrum-X BGP underlay template (and, in l3evpn, the EVPN overlay
RR-mesh, QoS PFC defaults, and the action-bound route-domain VRF template) plus
one variables-carrying profile per fabric switch, and reconcile each switch's
device customVariables (aggregates / is_evpn_rr) that the tenant VRF render reads.

Per-device variables (bgp_neighbors, /26 aggregates, RR selection, overlay
neighbors) are computed from the same topology+p2p plan as 'configure-switches';
the devices must already carry asn/loopbackAddress (run 'configure-switches'
first). The .j2 template bodies are supplied via the config's bgp.*TemplatePath
keys and uploaded as-is; the engine renders them server-side.

The configuration is supplied EITHER as a whole document via --config-source OR
built from the per-property flags below (the two are mutually exclusive). The
per-property flags cover the bgp section (--mode, --apply-mode, --template-path
/-label/-profile-priority, the --overlay-*, --pfc-*, --vrf-* template flags)
plus the plan sections it reads (--ordering, --topology-leaf-spine
[-links-per-pair], --topology-spine-super-spine[-links-per-pair],
--topology-leaf-host[...], --p2p-pool-* , --p2p-mtu).

Arguments:
  fabric_id    The ID or label of the fabric

Input (one of):
  --config-source   'pipe' or path to the YAML/JSON config (with a 'bgp' section).
  per-property flags as listed above and in --help.

Always available:
  --dry-run         Report the plan without writing.
  --verify-render   Render every device through the engine first; abort on any render error.

Examples:
  metalcloud-cli fabric configure-bgp 5 --config-source fabric-config.l3evpn.yaml --verify-render
  metalcloud-cli fabric configure-bgp 5 --mode l3evpn \
    --template-path ./freeform-bgp-underlay.j2 --overlay-template-path ./freeform-bgp-overlay.j2 \
    --pfc-template-path ./freeform-qos-pfc.j2 --vrf-template-path ./switch-configure-vrf-create.j2 \
    --topology-leaf-spine --topology-spine-super-spine --p2p-pool-leaf-spine 10.254.0.0/16 --dry-run

```
metalcloud-cli fabric configure-bgp fabric_id [flags]
```

### Options

```
      --apply-mode string                                  Profile apply mode: once | always (default once).
      --config-source string                               Source of the configuration (with a 'bgp' section). 'pipe' or path to a YAML/JSON file. Mutually exclusive with the per-property flags.
      --dry-run                                            Report the plan without making changes.
  -h, --help                                               help for configure-bgp
      --mode string                                        Fabric mode: purel3 | l3evpn (must match freeform.mode).
      --ordering string                                    Device ordering: managementAddress | identifierString | id. (default "managementAddress")
      --overlay-profile-priority int                       Overlay profile priority (default 61).
      --overlay-template-label string                      Overlay template label (default spectrumx-bgp-overlay).
      --overlay-template-path string                       Path to the EVPN overlay .j2 template body.
      --p2p-mtu int32                                      MTU applied to created links.
      --p2p-pool-leaf-host string                          Leaf->host /31 pool.
      --p2p-pool-leaf-spine string                         Leaf<->spine /31 pool.
      --p2p-pool-spine-super-spine string                  Spine<->superspine /31 pool.
      --pfc-profile-priority int                           PFC profile priority (default 62).
      --pfc-template-label string                          PFC template label (default spectrumx-qos-pfc).
      --pfc-template-path string                           Path to the QoS PFC .j2 template body.
      --profile-priority int                               Underlay profile priority (default 60).
      --template-label string                              Underlay template label (default spectrumx-bgp-underlay).
      --template-path string                               Path to the BGP underlay .j2 template body.
      --topology-leaf-host                                 Enable leaf->host downlinks (for the /26 aggregates).
      --topology-leaf-host-description-template string     Leaf->host description template.
      --topology-leaf-host-nic-names strings               Remote host NIC names (even count).
      --topology-leaf-host-node-count int                  Number of host port-pairs per leaf.
      --topology-leaf-host-nodes ints                      Exact 0-based node indices (mutually exclusive with node-count).
      --topology-leaf-host-port-pattern string             Leaf host port pattern, e.g. swp{port}s{sub}.
      --topology-leaf-spine                                Enable leaf<->spine pairing.
      --topology-leaf-spine-links-per-pair string          Leaf<->spine links per pair: 'auto' or an integer.
      --topology-spine-super-spine                         Enable spine<->superspine pairing (3-tier only).
      --topology-spine-super-spine-links-per-pair string   Spine<->superspine links per pair: 'auto' or an integer.
      --verify-render                                      Render each device through the engine before writing; abort on any render error.
      --vrf-template-label string                          VRF template label (default switch-configure-vrf-create).
      --vrf-template-path string                           Path to the route-domain VRF .j2 template body.
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API key
  -c, --config string          Config file path
  -d, --debug                  Set to enable debug logging
  -e, --endpoint string        MetalCloud API endpoint
  -f, --format string          Output format. Supported values are 'text','csv','md','json','yaml'. (default "text")
  -i, --insecure_skip_verify   Set to allow insecure transport
  -l, --log_file string        Log file path
  -v, --verbosity string       Log level verbosity (default "INFO")
```

### SEE ALSO

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

