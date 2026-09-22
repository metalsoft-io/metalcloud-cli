## metalcloud-cli point-to-point-link update

Update a point-to-point link

### Synopsis

Update the properties of a point-to-point link (label, name, description,
annotations) from a JSON/YAML configuration. The link's current revision is
sent as If-Match, so a concurrent change is rejected rather than overwritten.

Use 'update-config' for the staged configuration (MTU) and the
'allocation-strategy' / 'static-route' sub-commands for the config collections.

Required Arguments:
  link_id            The ID of the point-to-point link

Required Flags:
  --config-source    'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli point-to-point-link update-example > update.yaml
  metalcloud-cli point-to-point-link update 42 --config-source update.yaml
  echo '{"description":"leaf01 uplink"}' | metalcloud-cli p2p update 42 --config-source pipe

```
metalcloud-cli point-to-point-link update link_id [flags]
```

### Options

```
      --config-source string   Source of the updated link configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update
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

