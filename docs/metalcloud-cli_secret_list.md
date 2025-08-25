## metalcloud-cli secret list

List all secrets

### Synopsis

List all secrets in the current datacenter.

This command displays a table of all secrets with their basic information including:
- Secret ID and name
- Encrypted value (partial display for security)
- Usage type (credential, configuration, etc.)
- Owner information
- Creation and update timestamps

The output is formatted as a table by default and can be filtered or formatted
using global output flags.

Examples:
  # List all secrets
  metalcloud-cli secret list

  # List secrets with JSON output
  metalcloud-cli secret list --output json

  # List secrets with custom formatting
  metalcloud-cli secret list --output table --fields id,name,usage

```
metalcloud-cli secret list [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli secret](metalcloud-cli_secret.md)	 - Manage encrypted secrets for secure credential storage

