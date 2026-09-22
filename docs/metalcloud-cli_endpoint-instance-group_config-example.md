## metalcloud-cli endpoint-instance-group config-example

Print an example endpoint instance group create configuration

### Synopsis

Print an example configuration that can be edited and passed to
'endpoint-instance-group create --config-source'.

Examples:
  metalcloud-cli endpoint-instance-group config-example > group.json
  metalcloud-cli eig config-example -f yaml

```
metalcloud-cli endpoint-instance-group config-example [flags]
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

* [metalcloud-cli endpoint-instance-group](metalcloud-cli_endpoint-instance-group.md)	 - Endpoint instance group management

