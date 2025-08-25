## metalcloud-cli user ssh-key-delete

Delete an SSH key from a user account

### Synopsis

Remove an existing SSH key from a specific user account.

This command permanently deletes an SSH key from the user's account. Once deleted,
the key can no longer be used for authentication to instances.

Arguments:
  user_id                 The numeric ID of the user whose SSH key to delete
  key_id                  The numeric ID of the SSH key to delete

Examples:
  metalcloud-cli user ssh-key-delete 12345 67890
  metalcloud-cli user delete-ssh-key 12345 67890
  metalcloud-cli user remove-ssh-key 12345 67890

```
metalcloud-cli user ssh-key-delete user_id key_id [flags]
```

### Options

```
  -h, --help   help for ssh-key-delete
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

