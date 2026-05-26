## metalcloud-cli site device-auth-provider update

Update a device auth provider from JSON

```
metalcloud-cli site device-auth-provider update provider_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the device auth provider updates. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

