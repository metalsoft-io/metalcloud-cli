## metalcloud-cli endpoint-instance-group update-config

Update the configuration of an endpoint instance group

### Synopsis

Update the pending configuration of an endpoint instance group.

The current configuration is fetched first and its revision is sent as If-Match,
because writes below the /config path are guarded by the configuration revision.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance-group update-config 12 --config-source update.json
  echo '{"label":"new-label"}' | metalcloud-cli eig update-config 12 --config-source pipe

```
metalcloud-cli endpoint-instance-group update-config endpoint_instance_group_id [flags]
```

### Options

```
      --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli endpoint-instance-group](metalcloud-cli_endpoint-instance-group.md)	 - Endpoint instance group management

