## metalcloud-cli user ssh-key

Display a single SSH key of a user

### Synopsis

Display one SSH key of a specific user account.

Use 'user ssh-keys' to list the SSH keys of the user and obtain their IDs.

Required Arguments:
  user_id                 The numeric ID of the user owning the SSH key
  key_id                  The numeric ID of the SSH key to display

Examples:
  metalcloud-cli user ssh-key 12345 67890
  metalcloud-cli user get-ssh-key 12345 67890

```
metalcloud-cli user ssh-key user_id key_id [flags]
```

### Options

```
  -h, --help   help for ssh-key
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

* [metalcloud-cli user](metalcloud-cli_user.md)	 - Manage user accounts and their properties

