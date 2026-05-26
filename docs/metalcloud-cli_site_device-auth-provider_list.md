## metalcloud-cli site device-auth-provider list

List all device auth providers

```
metalcloud-cli site device-auth-provider list [flags...] [flags]
```

### Options

```
      --filter-kind strings      Filter providers by kind (e.g. tacacs).
      --filter-site-id strings   Filter providers by site ID.
      --filter-status strings    Filter providers by status (e.g. active, maintenance, disabled).
  -h, --help                     help for list
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

