## metalcloud-cli server-instance-group update-meta

Update the metadata of a server instance group

### Synopsis

Update the metadata (tags) of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Required Flags:
  --config-source string   Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance-group update-meta 1234 --config-source meta.json
  echo '{"tags":["prod"]}' | metalcloud-cli ig update-meta 1234 --config-source pipe

```
metalcloud-cli server-instance-group update-meta server_instance_group_id [flags]
```

### Options

```
      --config-source string   Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update-meta
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

* [metalcloud-cli server-instance-group](metalcloud-cli_server-instance-group.md)	 - Manage server instance groups within infrastructures

