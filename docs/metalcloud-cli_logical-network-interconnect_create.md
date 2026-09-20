## metalcloud-cli logical-network-interconnect create

Create a logical network interconnect

### Synopsis

Create a new logical network interconnect on top of an existing network fabric
interconnect.

Required Flags (one of):
  --config-source string        Source of the configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string                Label of the new logical network interconnect (used together with the flags below)

Required Flags (when not using --config-source):
  --name string                 Name of the new logical network interconnect
  --fabric-interconnect-id int  ID of the network fabric interconnect it runs on

Optional Flags (when not using --config-source):
  --kind string                 Interconnect kind (default "dci-evpn")

Examples:
  metalcloud logical-network-interconnect create --config-source lni.json
  metalcloud lni create --label dc1-dc2-ln --name "DC1 to DC2 network" --fabric-interconnect-id 12

```
metalcloud-cli logical-network-interconnect create [flags]
```

### Options

```
      --config-source string         Source of the new logical network interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.
      --fabric-interconnect-id int   ID of the network fabric interconnect it runs on.
  -h, --help                         help for create
      --kind string                  Kind of the logical network interconnect. (default "dci-evpn")
      --label string                 Label of the new logical network interconnect.
      --name string                  Name of the new logical network interconnect.
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

* [metalcloud-cli logical-network-interconnect](metalcloud-cli_logical-network-interconnect.md)	 - Logical network interconnect management

