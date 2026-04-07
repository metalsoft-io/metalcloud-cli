## metalcloud-cli network-device default-secrets

Manage network device default secrets

### Synopsis

Manage network device default secrets for network devices (switches).

Network device default secrets store secret values (such as passwords or keys) associated
with a specific network device identified by MAC address or serial number. These secrets
are encrypted and stored securely.

Examples:
  # List all network device default secrets
  metalcloud-cli network-device default-secrets list

  # Get specific secret information
  metalcloud-cli network-device default-secrets get 123

  # Create a new secret
  metalcloud-cli nd ds create --site-id 1 --mac-or-serial "AA:BB:CC:DD:EE:FF" --secret-name "admin_password" --secret-value "s3cur3"

### Options

```
  -h, --help   help for default-secrets
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

* [metalcloud-cli network-device](metalcloud-cli_network-device.md)	 - Manage network devices (switches) in the infrastructure
* [metalcloud-cli network-device default-secrets batch-create](metalcloud-cli_network-device_default-secrets_batch-create.md)	 - Create multiple network device default secrets from a CSV file
* [metalcloud-cli network-device default-secrets create](metalcloud-cli_network-device_default-secrets_create.md)	 - Create a new network device default secret
* [metalcloud-cli network-device default-secrets delete](metalcloud-cli_network-device_default-secrets_delete.md)	 - Delete a network device default secret
* [metalcloud-cli network-device default-secrets get](metalcloud-cli_network-device_default-secrets_get.md)	 - Get detailed information about a specific network device default secret
* [metalcloud-cli network-device default-secrets get-credentials](metalcloud-cli_network-device_default-secrets_get-credentials.md)	 - Retrieve the unencrypted secret value
* [metalcloud-cli network-device default-secrets list](metalcloud-cli_network-device_default-secrets_list.md)	 - List all network device default secrets
* [metalcloud-cli network-device default-secrets update](metalcloud-cli_network-device_default-secrets_update.md)	 - Update an existing network device default secret

