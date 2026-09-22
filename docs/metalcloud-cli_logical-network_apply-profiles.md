## metalcloud-cli logical-network apply-profiles

Apply a logical network profile to a logical network config

### Synopsis

Apply a logical network profile onto the config of an existing logical network.

The profile's allocation strategies replace the ones currently held by the
config. The config's own revision is sent as the If-Match entity tag.

Required Arguments:
  logical_network_id  The ID of the logical network
  profile_id          The ID of the logical network profile to apply

Examples:
  metalcloud-cli logical-network apply-profiles 12 3

```
metalcloud-cli logical-network apply-profiles logical_network_id profile_id [flags]
```

### Options

```
  -h, --help   help for apply-profiles
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

* [metalcloud-cli logical-network](metalcloud-cli_logical-network.md)	 - Manage logical networks within fabrics

