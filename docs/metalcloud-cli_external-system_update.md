## metalcloud-cli external-system update

Update an external system

### Synopsis

Update the label, name or annotations of an external system. Fields absent from
the configuration document are left unchanged.

The current revision is read first and sent as the If-Match entity tag, so the
update fails if somebody else changed the external system in the meantime.

Required Arguments:
  external_system_id       The numeric ID of the external system

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud external-system update 12 --config-source update.json
  echo '{"name":"New name"}' | metalcloud ext-system update 12 --config-source pipe

```
metalcloud-cli external-system update external_system_id [flags]
```

### Options

```
      --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli external-system](metalcloud-cli_external-system.md)	 - External system management

