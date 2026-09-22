## metalcloud-cli endpoint-instance-group network replace

Create or replace the network configuration of an endpoint instance group

### Synopsis

Create or replace (PUT) the network configuration of an endpoint instance group.

Without --config-source the operation is sent with no body, which creates the network
configuration if it does not exist yet and returns it otherwise. With --config-source
the supplied document is sent as the request body together with the If-Match of the
group configuration.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Optional Flags:
  --config-source string   Source of the network configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance-group network replace 12
  metalcloud-cli endpoint-instance-group network replace 12 --config-source networking.json

```
metalcloud-cli endpoint-instance-group network replace endpoint_instance_group_id [flags]
```

### Options

```
      --config-source string   Source of the network configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for replace
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

* [metalcloud-cli endpoint-instance-group network](metalcloud-cli_endpoint-instance-group_network.md)	 - Manage the network configuration of an endpoint instance group

