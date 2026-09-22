## metalcloud-cli endpoint-instance create

Create an endpoint instance in an infrastructure

### Synopsis

Create a new endpoint instance in an infrastructure.

The instance can be described either with a configuration file (--config-source) or
with individual flags (--endpoint-id, --label, ...).

Required Arguments:
  infrastructure_id_or_label   The ID or label of the infrastructure

Required Flags (one of):
  --config-source string   Source of the new endpoint instance configuration. Can be 'pipe' or path to a JSON/YAML file.
  --endpoint-id int        ID of the endpoint deployed by this instance

Optional Flags (when not using --config-source):
  --label string      Label of the new endpoint instance
  --group-id int      ID of the endpoint instance group the instance belongs to
  --tags strings      Tags of the new endpoint instance

Examples:
  metalcloud-cli endpoint-instance create 1234 --endpoint-id 7
  metalcloud-cli endpoint-instance create prod-env --config-source endpoint-instance.json
  cat endpoint-instance.yaml | metalcloud-cli ei new prod-env --config-source pipe

```
metalcloud-cli endpoint-instance create infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the new endpoint instance configuration. Can be 'pipe' or path to a JSON/YAML file.
      --endpoint-id int        ID of the endpoint deployed by this instance.
      --group-id int           ID of the endpoint instance group the instance belongs to.
  -h, --help                   help for create
      --label string           Label of the new endpoint instance.
      --tags strings           Tags of the new endpoint instance.
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

* [metalcloud-cli endpoint-instance](metalcloud-cli_endpoint-instance.md)	 - Endpoint instance management

