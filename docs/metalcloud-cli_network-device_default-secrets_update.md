## metalcloud-cli network-device default-secrets update

Update an existing network device default secret

### Synopsis

Update the secret value of an existing network device default secret.

Arguments:
  secrets_id    The ID of the network device default secrets to update (required)

Required Flags:
  --secret-value    New value of the secret

Examples:
  # Update the secret value
  metalcloud-cli network-device default-secrets update 123 --secret-value "new_s3cur3"

  # Using alias
  metalcloud-cli nd ds update 456 --secret-value "updated_password"

```
metalcloud-cli network-device default-secrets update <secrets_id> [flags]
```

### Options

```
  -h, --help                  help for update
      --secret-value string   New value of the secret
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

