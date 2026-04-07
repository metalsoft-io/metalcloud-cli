## metalcloud-cli network-device default-secrets delete

Delete a network device default secret

### Synopsis

Delete a network device default secret by ID.

This operation permanently removes the secret and cannot be undone.

Arguments:
  secrets_id    The ID of the network device default secrets to delete (required)

Examples:
  # Delete a secret
  metalcloud-cli network-device default-secrets delete 123

  # Using alias
  metalcloud-cli nd ds rm 456

```
metalcloud-cli network-device default-secrets delete <secrets_id> [flags]
```

### Options

```
  -h, --help   help for delete
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

