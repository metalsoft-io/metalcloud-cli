## metalcloud-cli site device-auth-provider create

Create a device auth provider from JSON

### Synopsis

Create a new device auth provider from a JSON document.

The configuration source must supply all required fields: label, name, siteId,
kind, ipAddress, port, sharedSecret, username. Run 'config-example' to print a
template.

```
metalcloud-cli site device-auth-provider create [flags]
```

### Options

```
      --config-source string   Source of the new device auth provider configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

* [metalcloud-cli site device-auth-provider](metalcloud-cli_site_device-auth-provider.md)	 - Manage device authentication providers (e.g. TACACS+)

