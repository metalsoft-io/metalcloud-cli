## metalcloud-cli server-default-credentials get

Get detailed information about specific server default credentials

### Synopsis

Get detailed information about specific server default credentials.

This command retrieves comprehensive information about a specific set of server default
credentials, including all metadata fields but not the actual password (use get-credentials
for that). The output includes ID, site ID, server serial number, MAC address, username,
and any optional metadata like rack information, inventory ID, and UUID.

Arguments:
  credentials_id    The ID of the server default credentials to retrieve (required)

Examples:
  # Get information about credentials with ID 123
  metalcloud-cli server-default-credentials get 123

  # Get credentials info using short alias
  metalcloud-cli sdc get 456

  # Get credentials info using alternate alias
  metalcloud-cli srv-dc get 789

```
metalcloud-cli server-default-credentials get <credentials_id> [flags]
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

* [metalcloud-cli server-default-credentials](metalcloud-cli_server-default-credentials.md)	 - Manage server default credentials and authentication settings

