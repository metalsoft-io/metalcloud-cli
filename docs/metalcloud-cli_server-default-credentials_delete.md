## metalcloud-cli server-default-credentials delete

Delete server default credentials

### Synopsis

Delete server default credentials by ID.

This command permanently removes a set of server default credentials from the system.
Once deleted, the credentials cannot be recovered and will no longer be available for
server provisioning or management operations.

Arguments:
  credentials_id    The ID of the server default credentials to delete (required)

Examples:
  # Delete credentials with ID 123
  metalcloud-cli server-default-credentials delete 123

  # Delete using short alias
  metalcloud-cli sdc rm 456

  # Delete using alternate alias
  metalcloud-cli srv-dc delete 789

```
metalcloud-cli server-default-credentials delete <credentials_id> [flags]
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

* [metalcloud-cli server-default-credentials](metalcloud-cli_server-default-credentials.md)	 - Manage server default credentials and authentication settings

