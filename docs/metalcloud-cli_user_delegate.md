## metalcloud-cli user delegate

Manage user delegates

### Synopsis

Manage the delegation relationships between user accounts.

A delegate user can act on the resources of the user that delegated access to them.
These commands allow you to:
- Grant delegate access to another user (add)
- Revoke delegate access (remove)
- List the users that delegated access to a user (parents)
- List the users a user delegated access to (children)

### Options

```
  -h, --help   help for delegate
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
* [metalcloud-cli user delegate add](metalcloud-cli_user_delegate_add.md)	 - Grant a user delegate access to another user
* [metalcloud-cli user delegate children](metalcloud-cli_user_delegate_children.md)	 - List the users a user delegated access to
* [metalcloud-cli user delegate parents](metalcloud-cli_user_delegate_parents.md)	 - List the users that delegated access to a user
* [metalcloud-cli user delegate remove](metalcloud-cli_user_delegate_remove.md)	 - Revoke the delegate access of a user

