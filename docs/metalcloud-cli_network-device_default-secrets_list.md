## metalcloud-cli network-device default-secrets list

List all network device default secrets

### Synopsis

List all network device default secrets with pagination support.

Flags:
  --page    Page number for paginated results (default: 0, which returns all results)
  --limit   Number of records per page, maximum 100 (default: 0, which returns all results)

Examples:
  # List all network device default secrets
  metalcloud-cli network-device default-secrets list

  # List with pagination
  metalcloud-cli nd ds list --page 2 --limit 10

```
metalcloud-cli network-device default-secrets list [flags]
```

### Options

```
  -h, --help        help for list
      --limit int   Number of records per page (max 100)
      --page int    Page number
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

