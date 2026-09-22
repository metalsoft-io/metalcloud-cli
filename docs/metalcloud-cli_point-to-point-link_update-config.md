## metalcloud-cli point-to-point-link update-config

Update the staged configuration of a point-to-point link

### Synopsis

Update the staged configuration of a point-to-point link (currently its MTU)
from a JSON/YAML configuration. The configuration object's own revision is sent
as If-Match - not the link's.

Required Arguments:
  link_id            The ID of the point-to-point link

Required Flags:
  --config-source    'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli point-to-point-link update-config-example > config.yaml
  metalcloud-cli point-to-point-link update-config 42 --config-source config.yaml
  echo '{"mtu":9216}' | metalcloud-cli p2p update-config 42 --config-source pipe

```
metalcloud-cli point-to-point-link update-config link_id [flags]
```

### Options

```
      --config-source string   Source of the updated link config. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update-config
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

