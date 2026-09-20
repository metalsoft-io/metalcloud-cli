## metalcloud-cli fabric bgp-session update

Update a BGP session of a fabric

### Synopsis

Update a BGP session of a network fabric from a JSON/YAML configuration.

Only the custom variables of a BGP session can be updated; the numbering, link
configuration and the link it runs on are fixed at creation time.

Required Arguments:
  fabric_id         The ID or name of the fabric
  bgp_session_id    The ID of the BGP session

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli fabric bgp-session update 12345 7 --config-source session.yaml
  echo '{"customVariables":{"asn":65000}}' | metalcloud-cli fabric bgp update 12345 7 --config-source pipe

```
metalcloud-cli fabric bgp-session update fabric_id bgp_session_id [flags]
```

### Options

```
      --config-source string   Source of the BGP session configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli fabric bgp-session](metalcloud-cli_fabric_bgp-session.md)	 - Manage the BGP sessions of a fabric

