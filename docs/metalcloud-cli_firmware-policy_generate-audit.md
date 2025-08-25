## metalcloud-cli firmware-policy generate-audit

Generate compliance audit report for a firmware policy

### Synopsis

Generate a compliance audit report for a firmware policy to analyze server firmware status.

This command analyzes the current firmware status of all servers that match the
specified policy rules and generates a detailed compliance report. The audit shows
which servers are compliant with the policy requirements, which need updates,
and provides detailed firmware version information.

Required arguments:
  policy_id               The unique identifier (numeric ID) of the firmware policy

The audit report includes:
  - List of servers matching the policy rules
  - Current firmware versions for each server component
  - Compliance status for each server
  - Recommended firmware updates
  - Servers that would be affected by policy execution

Examples:
  # Generate audit for firmware policy with ID 123
  metalcloud-cli firmware-policy generate-audit 123
  
  # Generate audit using alias
  metalcloud-cli fw-policy audit 456
  
  # Save audit results to file
  metalcloud-cli firmware-policy generate-audit 789 > audit-report.json

```
metalcloud-cli firmware-policy generate-audit policy_id [flags]
```

### Options

```
  -h, --help   help for generate-audit
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

