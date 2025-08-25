## metalcloud-cli user

Manage user accounts and their properties

### Synopsis

Comprehensive user management commands for creating, modifying, and managing user accounts.
These commands allow you to perform various operations on user accounts including:
- Creating individual or bulk users
- Managing user permissions and limits
- Handling SSH keys and authentication
- User lifecycle operations (archive/unarchive, suspend/unsuspend)
- Account management and configuration updates

All commands require appropriate permissions and most modification commands require
the user ID as a parameter. Use 'metalcloud-cli user list' to find user IDs.

### Options

```
  -h, --help   help for user
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli user archive](metalcloud-cli_user_archive.md)	 - Archive a user account to mark it as inactive
* [metalcloud-cli user change-account](metalcloud-cli_user_change-account.md)	 - Move a user to a different account
* [metalcloud-cli user config-update](metalcloud-cli_user_config-update.md)	 - Update comprehensive user configuration settings
* [metalcloud-cli user create](metalcloud-cli_user_create.md)	 - Create a new user account with specified properties
* [metalcloud-cli user create-bulk](metalcloud-cli_user_create-bulk.md)	 - Create multiple users in a single operation
* [metalcloud-cli user get](metalcloud-cli_user_get.md)	 - Display detailed information for a specific user
* [metalcloud-cli user limits](metalcloud-cli_user_limits.md)	 - Display resource limits for a specific user
* [metalcloud-cli user limits-update](metalcloud-cli_user_limits-update.md)	 - Update resource limits for a specific user
* [metalcloud-cli user list](metalcloud-cli_user_list.md)	 - List users with filtering and search options
* [metalcloud-cli user permissions](metalcloud-cli_user_permissions.md)	 - Display permissions for a specific user
* [metalcloud-cli user permissions-update](metalcloud-cli_user_permissions-update.md)	 - Update permissions for a specific user
* [metalcloud-cli user ssh-key-add](metalcloud-cli_user_ssh-key-add.md)	 - Add an SSH key to a user account
* [metalcloud-cli user ssh-key-delete](metalcloud-cli_user_ssh-key-delete.md)	 - Delete an SSH key from a user account
* [metalcloud-cli user ssh-keys](metalcloud-cli_user_ssh-keys.md)	 - Display SSH keys for a specific user
* [metalcloud-cli user suspend](metalcloud-cli_user_suspend.md)	 - Suspend a user account temporarily
* [metalcloud-cli user unarchive](metalcloud-cli_user_unarchive.md)	 - Unarchive a user account to restore access
* [metalcloud-cli user unsuspend](metalcloud-cli_user_unsuspend.md)	 - Unsuspend a user account to restore access

