## metalcloud-cli server-default-credentials config-example

Show a server default credentials update configuration example

### Synopsis

Show an example of the configuration accepted by 'server-default-credentials update'.

Examples:
  # Write the example to a file and edit it
  metalcloud-cli server-default-credentials config-example > credentials.json


```
metalcloud-cli server-default-credentials config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

* [metalcloud-cli server-default-credentials](metalcloud-cli_server-default-credentials.md)	 - Manage server default credentials and authentication settings

