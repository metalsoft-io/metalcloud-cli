## metalcloud-cli server-instance update-config

Update the configuration of a server instance

### Synopsis

Update the pending configuration of a server instance.

The current configuration is fetched first and its revision is sent as If-Match,
because writes below the /config path are guarded by the configuration revision.
Use 'server-instance config' to read the current configuration.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance update-config 5678 --config-source update.json
  echo '{"label":"new-label"}' | metalcloud-cli inst update-config 5678 --config-source pipe

```
metalcloud-cli server-instance update-config server_instance_id [flags]
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

