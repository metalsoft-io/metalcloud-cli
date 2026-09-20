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
* [metalcloud-cli user 2fa-disable](metalcloud-cli_user_2fa-disable.md)	 - Disable two-factor authentication
* [metalcloud-cli user 2fa-enable](metalcloud-cli_user_2fa-enable.md)	 - Enable two-factor authentication
* [metalcloud-cli user 2fa-generate-secret](metalcloud-cli_user_2fa-generate-secret.md)	 - Generate a new 2FA secret for setting up an authenticator app
* [metalcloud-cli user api-key](metalcloud-cli_user_api-key.md)	 - Get the current user's API key
* [metalcloud-cli user api-key-regenerate](metalcloud-cli_user_api-key-regenerate.md)	 - Regenerate the current user's API key
* [metalcloud-cli user archive](metalcloud-cli_user_archive.md)	 - Archive a user account to mark it as inactive
* [metalcloud-cli user change-account](metalcloud-cli_user_change-account.md)	 - Move a user to a different account
* [metalcloud-cli user change-password](metalcloud-cli_user_change-password.md)	 - Change the password of the current user
* [metalcloud-cli user config](metalcloud-cli_user_config.md)	 - Display the configuration of a user
* [metalcloud-cli user config-update](metalcloud-cli_user_config-update.md)	 - Update comprehensive user configuration settings
* [metalcloud-cli user create](metalcloud-cli_user_create.md)	 - Create a new user account with specified properties
* [metalcloud-cli user create-bulk](metalcloud-cli_user_create-bulk.md)	 - Create multiple users in a single operation
* [metalcloud-cli user delegate](metalcloud-cli_user_delegate.md)	 - Manage user delegates
* [metalcloud-cli user delete](metalcloud-cli_user_delete.md)	 - Delete a user and erase their personal information
* [metalcloud-cli user get](metalcloud-cli_user_get.md)	 - Display detailed information for a specific user
* [metalcloud-cli user initiate-email-change](metalcloud-cli_user_initiate-email-change.md)	 - Start changing the e-mail address of the current user
* [metalcloud-cli user initiate-password-reset](metalcloud-cli_user_initiate-password-reset.md)	 - Send a password reset message to an e-mail address
* [metalcloud-cli user limits](metalcloud-cli_user_limits.md)	 - Display resource limits for a specific user
* [metalcloud-cli user list](metalcloud-cli_user_list.md)	 - List users with filtering and search options
* [metalcloud-cli user permissions](metalcloud-cli_user_permissions.md)	 - List the permissions of the current user
* [metalcloud-cli user regenerate-jwt-salt](metalcloud-cli_user_regenerate-jwt-salt.md)	 - Regenerate the JWT salt of the current user
* [metalcloud-cli user resend-email-verification](metalcloud-cli_user_resend-email-verification.md)	 - Resend the e-mail verification message to a user
* [metalcloud-cli user resend-invitation](metalcloud-cli_user_resend-invitation.md)	 - Resend the platform invitation to a user
* [metalcloud-cli user reset-password](metalcloud-cli_user_reset-password.md)	 - Consume a password reset token
* [metalcloud-cli user send-password-reset](metalcloud-cli_user_send-password-reset.md)	 - Send a password reset message to a user (admin)
* [metalcloud-cli user set-password](metalcloud-cli_user_set-password.md)	 - Set the password for a user (admin only)
* [metalcloud-cli user ssh-key](metalcloud-cli_user_ssh-key.md)	 - Display a single SSH key of a user
* [metalcloud-cli user ssh-key-add](metalcloud-cli_user_ssh-key-add.md)	 - Add an SSH key to a user account
* [metalcloud-cli user ssh-key-delete](metalcloud-cli_user_ssh-key-delete.md)	 - Delete an SSH key from a user account
* [metalcloud-cli user ssh-keys](metalcloud-cli_user_ssh-keys.md)	 - Display SSH keys for a specific user
* [metalcloud-cli user suspend](metalcloud-cli_user_suspend.md)	 - Suspend a user account temporarily
* [metalcloud-cli user suspend-reasons](metalcloud-cli_user_suspend-reasons.md)	 - List the suspend reasons recorded for a user
* [metalcloud-cli user unarchive](metalcloud-cli_user_unarchive.md)	 - Unarchive a user account to restore access
* [metalcloud-cli user unsuspend](metalcloud-cli_user_unsuspend.md)	 - Unsuspend a user account to restore access
* [metalcloud-cli user update-meta](metalcloud-cli_user_update-meta.md)	 - Update the metadata of a user
* [metalcloud-cli user verify-email](metalcloud-cli_user_verify-email.md)	 - Verify an e-mail address with a verification token

