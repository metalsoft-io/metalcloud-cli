## metalcloud-cli server-instance-group network acl config-example

Print an example security rule configuration

### Synopsis

Print an example configuration that can be edited and passed to
'server-instance-group network acl add --config-source'.

Examples:
  metalcloud-cli server-instance-group network acl config-example > rule.json
  metalcloud-cli ig net acl config-example -f yaml

```
metalcloud-cli server-instance-group network acl config-example [flags]
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

* [metalcloud-cli server-instance-group network acl](metalcloud-cli_server-instance-group_network_acl.md)	 - Manage the security rules of a network connection

