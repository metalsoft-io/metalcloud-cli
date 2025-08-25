## metalcloud-cli secret get

Get secret details by ID

### Synopsis

Get detailed information about a specific secret by its ID.

This command retrieves and displays comprehensive information about a secret,
including its name, encrypted value (partial display for security), usage type,
owner details, and timestamps. The secret ID must be provided as a numeric value.

Arguments:
  secret_id          Numeric ID of the secret to retrieve (required)

Examples:
  # Get details of secret with ID 123
  metalcloud-cli secret get 123

  # Get secret details with JSON output
  metalcloud-cli secret get 456 --output json

  # Get secret details with custom field selection
  metalcloud-cli secret get 789 --output table --fields id,name,usage,createdTimestamp

```
metalcloud-cli secret get secret_id [flags]
```

### Options

```
  -h, --help   help for get
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

