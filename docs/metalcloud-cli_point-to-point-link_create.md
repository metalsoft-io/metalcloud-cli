## metalcloud-cli point-to-point-link create

Create a point-to-point link

### Synopsis

Create a point-to-point link from a JSON/YAML configuration.

The configuration may include:
- interfaceA / interfaceB: { type: "network_equipment_interface", interfaceId: <id> }
  (omit interfaceB for a half-connected link)
- description, mtu, routingActivation ("default" or "while_transporting_logical_network")
- ipv4.subnetAllocationStrategies: one or more strategies staged on create
  (e.g. a manual strategy with a subnetId and interfaceABinding)

Required Flags:
  --config-source    'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli point-to-point-link config-example > link.json
  metalcloud-cli point-to-point-link create --config-source link.json
  cat link.json | metalcloud-cli p2p create --config-source pipe

```
metalcloud-cli point-to-point-link create [flags]
```

### Options

```
      --config-source string   Source of the new link configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for create
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

* [metalcloud-cli point-to-point-link](metalcloud-cli_point-to-point-link.md)	 - Manage point-to-point links between network interfaces

