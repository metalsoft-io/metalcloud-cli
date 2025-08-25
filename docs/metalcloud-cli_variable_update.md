## metalcloud-cli variable update

Update an existing variable

### Synopsis

Update an existing variable with new configuration values.

This command updates a variable's configuration using data provided through a JSON file
or piped input. You can modify the variable's name, value, and usage type.

Required Arguments:
  variable_id    Numeric ID of the variable to update

Required Flags:
  --config-source    Source of the variable configuration updates

Configuration Source Options:
  pipe              Read configuration from stdin (use with echo or cat)
  /path/to/file     Read configuration from specified JSON file

Configuration Format:
  The configuration must be valid JSON. You can include any of these fields:
  {
    "name": "new-variable-name",       // Optional: New variable name
    "value": {                         // Optional: New key-value pairs
      "key1": "new-value1",
      "key2": "new-value2"
    },
    "usage": "specific"                // Optional: New usage type
  }

  Note: Only the fields you specify will be updated. Omitted fields remain unchanged.

Required Permissions:
  VARIABLES_AND_SECRETS_WRITE

Examples:
  # Update variable from JSON file
  metalcloud-cli variable update 123 --config-source /path/to/update.json
  
  # Update variable from stdin using pipe
  echo '{"name":"updated-name","value":{"new-key":"new-value"}}' | metalcloud-cli variable update 123 --config-source pipe
  
  # Update only the value using cat
  cat value-update.json | metalcloud-cli variable update 456 --config-source pipe
  
  # Using alias
  metalcloud-cli var edit 789 --config-source /tmp/variable-update.json

```
metalcloud-cli variable update variable_id [flags]
```

### Options

```
      --config-source string   Source of the variable configuration updates. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

* [metalcloud-cli variable](metalcloud-cli_variable.md)	 - Manage variables for infrastructure configuration

