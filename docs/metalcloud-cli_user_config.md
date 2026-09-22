## metalcloud-cli user config

Display the configuration of a user

### Synopsis

Display the configuration object of a specific user account.

The configuration holds the settings written by 'user config-update': display name,
access level, language, brand, login state and password policy flags.

Required Arguments:
  user_id                 The numeric ID of the user whose configuration to display

Examples:
  metalcloud-cli user config 12345
  metalcloud-cli user get-config 12345

```
metalcloud-cli user config user_id [flags]
```

### Options

```
  -h, --help   help for config
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

