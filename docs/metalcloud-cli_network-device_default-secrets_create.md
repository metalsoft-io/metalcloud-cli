## metalcloud-cli network-device default-secrets create

Create a new network device default secret

### Synopsis

Create a new network device default secret.

Required Flags:
  --site-id          Site ID where the network device is located
  --mac-or-serial    MAC address or serial number of the network device
  --secret-name      Name of the secret
  --secret-value     Value of the secret

Examples:
  # Create a new secret
  metalcloud-cli network-device default-secrets create \
    --site-id 1 \
    --mac-or-serial "AA:BB:CC:DD:EE:FF" \
    --secret-name "admin_password" \
    --secret-value "s3cur3"

  # Using alias
  metalcloud-cli nd ds create --site-id 2 --mac-or-serial "SN123456" --secret-name "enable_secret" --secret-value "mypass"

```
metalcloud-cli network-device default-secrets create [flags]
```

### Options

```
  -h, --help                   help for create
      --mac-or-serial string   MAC address or serial number of the network device
      --secret-name string     Name of the secret
      --secret-value string    Value of the secret
      --site-id float32        Site ID
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

