## metalcloud-cli firmware-policy get

Get detailed information about a specific firmware policy

### Synopsis

Get detailed information about a specific firmware policy including its configuration,
rules, associated server instance groups, and current status.

This command retrieves and displays all available information for a single firmware
policy, including the rules that determine which servers the policy applies to,
the firmware upgrade action to be performed, and any server instance groups
that are linked to this policy.

Required arguments:
  policy_id               The unique identifier (numeric ID) of the firmware policy

Examples:
  # Get details for firmware policy with ID 123
  metalcloud-cli firmware-policy get 123
  
  # Show policy details using alias
  metalcloud-cli fw-policy show 456

```
metalcloud-cli firmware-policy get policy_id [flags]
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

* [metalcloud-cli firmware-policy](metalcloud-cli_firmware-policy.md)	 - Manage server firmware upgrade policies and global firmware configurations

