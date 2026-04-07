## metalcloud-cli network-device default-secrets get

Get detailed information about a specific network device default secret

### Synopsis

Get detailed information about a specific network device default secret.

This returns metadata about the secret (ID, site, MAC/serial, name, timestamps)
but not the actual secret value. Use get-credentials to retrieve the secret value.

Arguments:
  secrets_id    The ID of the network device default secrets to retrieve (required)

Examples:
  # Get secret information
  metalcloud-cli network-device default-secrets get 123

  # Using alias
  metalcloud-cli nd ds show 456

```
metalcloud-cli network-device default-secrets get <secrets_id> [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli network-device default-secrets](metalcloud-cli_network-device_default-secrets.md)	 - Manage network device default secrets

