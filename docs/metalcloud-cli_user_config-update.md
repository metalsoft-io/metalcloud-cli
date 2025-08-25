## metalcloud-cli user config-update

Update comprehensive user configuration settings

### Synopsis

Update comprehensive configuration settings for a specific user account.

This command allows updating various user properties including display name, email, access level,
and other account settings. The configuration is provided through a JSON file or pipe.

Arguments:
  user_id                 The numeric ID of the user whose configuration to update

Required Flags:
  --config-source         Source of user configuration (JSON file path or 'pipe')

Configuration File Format (JSON):
  {
    "displayName": "Updated Name",
    "accessLevel": "admin",
    "emailVerified": true,
    "language": "en"
  }

Examples:
  metalcloud-cli user config-update 12345 --config-source config.json
  echo '{"displayName": "New Name", "accessLevel": "admin"}' | metalcloud-cli user config-update 12345 --config-source pipe

```
metalcloud-cli user config-update user_id [flags]
```

### Options

```
      --config-source string   Source of the user configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for config-update
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

