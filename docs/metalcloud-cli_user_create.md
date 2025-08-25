## metalcloud-cli user create

Create a new user account with specified properties

### Synopsis

Create a new user account in the system with comprehensive configuration options.

This command allows creating users either through individual command-line flags or by providing
a JSON configuration file/pipe. The user can be associated with an existing account or a new
account can be created automatically.

Required Flags (when not using --config-source):
  --email                 User's email address (required, used as login)
  --password              User's password (required for CLI creation)

Optional Flags:
  --config-source         Source of user configuration (JSON file path or 'pipe')
  --display-name          User's display name (defaults to email if not provided)
  --access-level          User access level: admin, user, readonly (default: user)
  --email-verified        Mark user email as verified (default: false)
  --account-id            Associate user with existing account ID
  --create-with-account   Create a new account for the user (mutually exclusive with --account-id)

Dependencies:
  - --email and --password are required together when not using --config-source
  - --config-source is mutually exclusive with --email
  - --account-id and --create-with-account are mutually exclusive

Configuration File Format (JSON):
  {
    "displayName": "John Doe",
    "email": "john.doe@company.com",
    "password": "securePassword123",
    "accessLevel": "user",
    "emailVerified": true,
    "accountId": 12345
  }

```
metalcloud-cli user create [flags]
```

### Examples

```
  # Create user with command-line flags
  metalcloud-cli user create --email test.user@metalsoft.io --password secret --access-level user
  
  # Create user with additional properties
  metalcloud-cli user create --email test.user@metalsoft.io --password secret --access-level user --display-name "Test User" --email-verified true --account-id 12345
  
  # Create user with new account
  metalcloud-cli user create --email admin@company.com --password admin123 --access-level admin --create-with-account
  
  # Create user from JSON file
  metalcloud-cli user create --config-source user1.json
  
  # Create user from pipe
  echo '{"email": "test.user@metalsoft.io", "password": "secret", "accessLevel": "user", "displayName": "Test User"}' | metalcloud-cli user create --config-source pipe
```

### Options

```
      --access-level string    Access level (e.g., 'admin', 'user')
      --account-id int         Account ID to associate the user with
      --config-source string   Source of the new user configuration. Can be 'pipe' or path to a JSON file.
      --create-with-account    Create new account for the user
      --display-name string    User's display name
      --email string           User's email address
      --email-verified         Set the user email as verified
  -h, --help                   help for create
      --password string        User's password (if not provided, a random password will be generated)
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

