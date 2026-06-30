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

Arguments:
  fabric_id    The ID or label of the fabric

Required Flags:
  --config-source   'pipe' or path to the YAML/JSON config (with a 'bgp' section).

Optional Flags:
  --dry-run         Report the plan without writing.
  --verify-render   Render every device through the engine first; abort on any render error.

Examples:
  metalcloud-cli fabric configure-bgp 5 --config-source fabric-config.l3evpn.yaml --verify-render
  metalcloud-cli fabric configure-bgp my-fabric --config-source fabric-config.yaml --dry-run

```
metalcloud-cli fabric configure-bgp fabric_id [flags]
```

### Options

```
      --config-source string   Source of the configuration (with a 'bgp' section). 'pipe' or path to a YAML/JSON file.
      --dry-run                Report the plan without making changes.
  -h, --help                   help for configure-bgp
      --verify-render          Render each device through the engine before writing; abort on any render error.
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

