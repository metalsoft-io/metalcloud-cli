## metalcloud-cli network-fabric-interconnect create

Create a network fabric interconnect

### Synopsis

Create a new network fabric interconnect.

The interconnect can be described either with a configuration file (--config-source)
or with individual flags (--label, --bgp-template-id, ...).

Required Flags (one of):
  --config-source string     Source of the interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string             Label of the new interconnect (used together with the flags below)

Optional Flags (when not using --config-source):
  --bgp-template-id int      ID of the BGP interconnect configuration template (required with --label)
  --type string              Interconnect type (default "dci-evpn")
  --name string              Display name
  --description string       Description
  --transport-id int         ID of the transport used by the interconnect

Examples:
  metalcloud network-fabric-interconnect create --config-source interconnect.json
  cat interconnect.yaml | metalcloud nfi create --config-source pipe
  metalcloud nfi create --label dc1-dc2 --bgp-template-id 3 --name "DC1 to DC2"

```
metalcloud-cli network-fabric-interconnect create [flags]
```

### Options

```
      --bgp-template-id int    ID of the BGP interconnect configuration template.
      --config-source string   Source of the new interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.
      --description string     Description of the new interconnect.
  -h, --help                   help for create
      --label string           Label of the new interconnect.
      --name string            Display name of the new interconnect.
      --transport-id int       ID of the transport used by the interconnect.
      --type string            Interconnect type. (default "dci-evpn")
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

* [metalcloud-cli network-fabric-interconnect](metalcloud-cli_network-fabric-interconnect.md)	 - Network fabric interconnect management

