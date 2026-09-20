## metalcloud-cli dhcp-reservation config-example

Print example DHCP reservation configurations

### Synopsis

Print one example create body per allocation kind (manual and auto).

Edit one of them and pass it to 'create' via --config-source.

Examples:
  metalcloud-cli dhcp-reservation config-example
  metalcloud-cli dhcp-reservation config-example -f yaml > reservation.yaml

```
metalcloud-cli dhcp-reservation config-example [flags]
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

* [metalcloud-cli dhcp-reservation](metalcloud-cli_dhcp-reservation.md)	 - Manage site DHCP reservations

