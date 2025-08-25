## metalcloud-cli account create

Create a new account with custom configuration

### Synopsis

Create a new account in the MetalCloud platform with custom configuration.

This command creates a new account using configuration data provided through 
a JSON file or standard input. The configuration must include all required 
account properties such as name, description, and any custom settings.

Required Permissions:
  - users:write

Required Flags:
  --config-source    Source of the new account configuration

Flag Details:
  --config-source string    Source of the new account configuration. 
                           Can be 'pipe' to read from stdin or path to a JSON file.
                           The JSON should contain account properties like name, 
                           description, and other account settings.

Examples:
  # Create account from JSON file
  metalcloud-cli account create --config-source /path/to/account.json

  # Create account from stdin
  echo '{"name":"test-account","description":"Test account"}' | metalcloud-cli account create --config-source pipe

  # Using alias
  metalcloud-cli account new --config-source account-config.json

```
metalcloud-cli account create [flags]
```

### Options

```
      --config-source string   Source of the new account configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

* [metalcloud-cli account](metalcloud-cli_account.md)	 - Manage user accounts and account-related operations

