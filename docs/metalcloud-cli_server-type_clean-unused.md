## metalcloud-cli server-type clean-unused

Remove the server types that are no longer used

### Synopsis

Remove every server type that no server references anymore.

Server types are created automatically when servers with a new hardware
configuration are registered; this command cleans up the ones left behind when
those servers are removed.

Examples:
  # Remove all unused server types
  metalcloud-cli server-type clean-unused


```
metalcloud-cli server-type clean-unused [flags]
```

### Options

```
  -h, --help   help for clean-unused
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

* [metalcloud-cli server-type](metalcloud-cli_server-type.md)	 - Manage server types and hardware configurations

