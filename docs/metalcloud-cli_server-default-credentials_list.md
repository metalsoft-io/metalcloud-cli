## metalcloud-cli server-default-credentials list

List all server default credentials

### Synopsis

List all server default credentials with pagination support.

This command displays server default credentials in a tabular format showing ID, site ID,
server serial number, MAC address, username, and optional metadata like rack information.

Flags:
  --page    Page number for paginated results (default: 0, which returns all results)
  --limit   Number of records per page, maximum 100 (default: 0, which returns all results)

Examples:
  # List all server default credentials
  metalcloud-cli server-default-credentials list

  # List credentials with pagination (page 2, 10 records per page)
  metalcloud-cli sdc list --page 2 --limit 10

  # List first 25 credentials
  metalcloud-cli srv-dc ls --limit 25

```
metalcloud-cli server-default-credentials list [flags]
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

* [metalcloud-cli server-default-credentials](metalcloud-cli_server-default-credentials.md)	 - Manage server default credentials and authentication settings

