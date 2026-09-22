## metalcloud-cli server-instance update-interface-config

Update the configuration of a server instance interface

### Synopsis

Update the pending configuration of one network interface of a server instance.

The interface is fetched first and the revision of its configuration is sent as
If-Match, because writes below the /config path are guarded by the configuration
revision.

Required Arguments:
  server_instance_id  The numeric ID of the server instance
  interface_id        The numeric ID of the interface

Required Flags:
  --config-source string   Source of the updated interface configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance update-interface-config 5678 7 --config-source iface.json
  echo '{"networkId":42}' | metalcloud-cli inst update-iface 5678 7 --config-source pipe

```
metalcloud-cli server-instance update-interface-config server_instance_id interface_id [flags]
```

### Options

```
      --config-source string   Source of the updated interface configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update-interface-config
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

