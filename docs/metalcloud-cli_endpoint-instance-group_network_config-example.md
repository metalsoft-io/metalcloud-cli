## metalcloud-cli endpoint-instance-group network config-example

Print an example network connection configuration

### Synopsis

Print an example configuration that can be edited and passed to
'endpoint-instance-group network connect --config-source'.

Examples:
  metalcloud-cli endpoint-instance-group network config-example > connection.json
  metalcloud-cli eig net config-example -f yaml

```
metalcloud-cli endpoint-instance-group network config-example [flags]
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

* [metalcloud-cli endpoint-instance-group network](metalcloud-cli_endpoint-instance-group_network.md)	 - Manage the network configuration of an endpoint instance group

