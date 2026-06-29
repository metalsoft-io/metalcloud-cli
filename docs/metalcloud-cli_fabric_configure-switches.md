## metalcloud-cli fabric configure-switches

Configure all switches of a fabric from a declarative YAML/JSON

### Synopsis

Configure every network device attached to a fabric from one declarative
configuration: hostnames (identifierString), ASNs, loopback IPs, physical-port
enable + interface descriptions, and point-to-point links with deterministic
/31 IPAM subnets.

Each feature section is optional - omit one to skip that step. Every step is
idempotent: current state is read first and only differences are written. Use
--dry-run to compute and preview the full plan without making any changes.

Arguments:
  fabric_id    The ID or label of the fabric to configure

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a YAML/JSON config file.

Optional Flags:
  --dry-run         Compute the plan and report what would change, without writing.

Config sections: hostname, asn, loopback, topology (leafSpine/spineSuperSpine/
leafHost), p2p, descriptionTemplate, enablePhysicalPorts, ordering.

Examples:
  metalcloud-cli fabric configure-switches 5 --config-source fabric-config.yaml --dry-run
  metalcloud-cli fabric configure-switches my-fabric --config-source fabric-config.yaml
  cat fabric-config.yaml | metalcloud-cli fabric configure-switches 5 --config-source pipe

```
metalcloud-cli fabric configure-switches fabric_id [flags]
```

### Options

```
      --config-source string   Source of the switch configuration. Can be 'pipe' or path to a YAML/JSON file.
      --dry-run                Compute and preview the plan without making any changes.
  -h, --help                   help for configure-switches
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

