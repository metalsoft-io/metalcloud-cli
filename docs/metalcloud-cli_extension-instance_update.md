## metalcloud-cli extension-instance update

Modify existing extension instance configuration

### Synopsis

Modify existing extension instance configuration with updated parameters.

This command allows you to update the configuration of an existing extension instance.
The updated configuration must be provided through the --config-source flag, which
accepts either 'pipe' for stdin input or a path to a JSON file containing the
updated configuration.

Use this command to modify input variables, change labels, or update other
configurable parameters of a deployed extension instance. The instance will
be reconfigured with the new settings while maintaining its deployment state.

Arguments:
  extension_instance_id    The unique ID of the extension instance to update

Required Flags:
  --config-source string   Source of the updated configuration (required)
                          Can be 'pipe' for stdin or path to a JSON file

JSON Configuration Format:
  {
    "label": "updated-instance-label",
    "inputVariables": [
      {"label": "variable1", "value": "new-value1"},
      {"label": "variable2", "value": "new-value2"}
    ]
  }

Examples:
  # Update from JSON file
  metalcloud extension-instance update 12345 --config-source ./updated-config.json
  
  # Update from pipe
  echo '{"label": "new-label"}' | metalcloud extension-instance update 12345 --config-source pipe
  
  # Update input variables
  metalcloud extension-instance update 12345 --config-source ./new-variables.json
  
  # Edit with alias
  metalcloud ext-inst edit 67890 --config-source updated-config.json

```
metalcloud-cli extension-instance update extension_instance_id [flags]
```

### Options

```
      --config-source string   Source of the extension instance configuration updates. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli extension-instance](metalcloud-cli_extension-instance.md)	 - Manage extension instances within infrastructure deployments

