## metalcloud-cli fabric configure-freeform

Register the base freeform template + per-switch profiles (step 8a)

### Synopsis

Register the Spectrum-X base freeform device-configuration template (hostname,
RoCE/QoS, adaptive routing, telemetry, BFD; in l3evpn the EVPN/VXLAN data plane)
and one variables-carrying profile per fabric switch, idempotently.

Per-device variables (mode, hgx_prefix, and for l3evpn leaves nve_source) are
computed from the same topology/loopback plan as 'configure-switches'. The .j2
template body is supplied via the config's freeform.templatePath and uploaded
as-is; the engine renders it server-side.

The configuration is supplied EITHER as a whole document via --config-source OR
built from the per-property flags below (the two are mutually exclusive). The
per-property flags cover the freeform section (--mode, --template-path,
--template-label, --profile-priority, --apply-mode, --hgx-prefix) plus the plan
sections it reads (--ordering, --topology-leaf-spine[-links-per-pair],
--topology-leaf-host[-node-count|...], --p2p-pool-leaf-host, ...).

Arguments:
  fabric_id    The ID or label of the fabric

Input (one of):
  --config-source   'pipe' or path to the YAML/JSON config (with a 'freeform' section).
  per-property flags as listed above and in --help.

Always available:
  --dry-run         Report the plan without writing.
  --verify-render   Render every device through the engine first; abort on any render error.

Examples:
  metalcloud-cli fabric configure-freeform 5 --config-source fabric-config.l3evpn.yaml --verify-render
  metalcloud-cli fabric configure-freeform 5 --mode l3evpn --template-path ./freeform-device-config.j2 \
    --topology-leaf-spine --dry-run

```
metalcloud-cli fabric configure-freeform fabric_id [flags]
```

### Options

```
      --apply-mode string                                  Profile apply mode: once | always (default once).
      --config-source string                               Source of the configuration (with a 'freeform' section). 'pipe' or path to a YAML/JSON file. Mutually exclusive with the per-property flags.
      --dry-run                                            Report the plan without making changes.
  -h, --help                                               help for configure-freeform
      --hgx-prefix string                                  Tenant HGX supernet prefix (default: derived from the leafHost pool).
      --mode string                                        Fabric mode: purel3 | l3evpn (must match bgp.mode).
      --ordering string                                    Device ordering: managementAddress | identifierString | id. (default "managementAddress")
      --p2p-mtu int32                                      MTU applied to created links.
      --p2p-pool-leaf-host string                          Leaf->host /31 pool.
      --p2p-pool-leaf-spine string                         Leaf<->spine /31 pool.
      --p2p-pool-spine-super-spine string                  Spine<->superspine /31 pool.
      --profile-priority int                               Profile priority (default 50).
      --template-label string                              Find-or-create template label (default spectrumx-freeform).
      --template-path string                               Path to the base freeform .j2 template body.
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

