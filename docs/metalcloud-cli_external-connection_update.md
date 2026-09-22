## metalcloud-cli external-connection update

Update an external connection

### Synopsis

Update the label or name of an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud external-connection update 12 --config-source update.json
  echo '{"name":"new name"}' | metalcloud ext-conn update dc1-ext --config-source pipe

```
metalcloud-cli external-connection update external_connection_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the updated external connection configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli external-connection](metalcloud-cli_external-connection.md)	 - External connection management

