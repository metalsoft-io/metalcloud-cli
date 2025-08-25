## metalcloud-cli server-default-credentials get-credentials

Retrieve unencrypted password for server default credentials

### Synopsis

Retrieve the unencrypted password for specific server default credentials.

This command returns the decrypted username and password for a specific set of server
default credentials. Use this when you need the actual password values for authentication
or configuration purposes. The password is decrypted server-side and transmitted securely.

Arguments:
  credentials_id    The ID of the server default credentials to retrieve password for (required)

Examples:
  # Get password for credentials with ID 123
  metalcloud-cli server-default-credentials get-credentials 123

  # Get password using alias
  metalcloud-cli sdc get-password 456

  # Get password using short alias
  metalcloud-cli srv-dc password 789

```
metalcloud-cli server-default-credentials get-credentials <credentials_id> [flags]
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

* [metalcloud-cli server-default-credentials](metalcloud-cli_server-default-credentials.md)	 - Manage server default credentials and authentication settings

