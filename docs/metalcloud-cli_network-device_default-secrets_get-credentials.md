## metalcloud-cli network-device default-secrets get-credentials

Retrieve the unencrypted secret value

### Synopsis

Retrieve the unencrypted secret value for a specific network device default secret.

The secret value is decrypted server-side and returned in plain text.

Arguments:
  secrets_id    The ID of the network device default secrets (required)

Examples:
  # Get the secret value
  metalcloud-cli network-device default-secrets get-credentials 123

  # Using alias
  metalcloud-cli nd ds get-secret 456

```
metalcloud-cli network-device default-secrets get-credentials <secrets_id> [flags]
```

### Options

```
  -h, --help   help for get-credentials
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

