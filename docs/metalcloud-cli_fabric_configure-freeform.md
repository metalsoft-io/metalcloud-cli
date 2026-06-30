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

Arguments:
  fabric_id    The ID or label of the fabric

Required Flags:
  --config-source   'pipe' or path to the YAML/JSON config (with a 'freeform' section).

Optional Flags:
  --dry-run         Report the plan without writing.
  --verify-render   Render every device through the engine first; abort on any render error.

Examples:
  metalcloud-cli fabric configure-freeform 5 --config-source fabric-config.l3evpn.yaml --verify-render
  metalcloud-cli fabric configure-freeform my-fabric --config-source fabric-config.yaml --dry-run

```
metalcloud-cli fabric configure-freeform fabric_id [flags]
```

### Options

```
      --config-source string   Source of the configuration (with a 'freeform' section). 'pipe' or path to a YAML/JSON file.
      --dry-run                Report the plan without making changes.
  -h, --help                   help for configure-freeform
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

