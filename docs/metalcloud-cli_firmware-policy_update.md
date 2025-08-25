## metalcloud-cli firmware-policy update

Update an existing firmware upgrade policy

### Synopsis

Update an existing firmware upgrade policy with new configuration.

This command allows you to modify an existing firmware policy by providing
updated configuration data. You can change the policy's label, action, rules,
and associated server instance groups. The policy ID cannot be changed.

Required arguments:
  policy_id               The unique identifier (numeric ID) of the firmware policy

Required flags:
  --config-source         Source of the firmware policy configuration updates
                          Values: 'pipe' (read from stdin) or path to JSON file

The configuration JSON can include any of these fields:
  - label: (optional) Updated descriptive name for the policy
  - action: (optional) Updated upgrade action (e.g., "upgrade", "downgrade") 
  - rules: (optional) Updated array of rules defining server selection criteria
  - userIdOwner: (optional) Updated user ID of the policy owner
  - serverInstanceGroupIds: (optional) Updated array of server instance group IDs

Note: Only provide the fields you want to update. Missing fields will retain
their current values.

Examples:
  # Update policy from JSON file
  metalcloud-cli firmware-policy update 123 --config-source policy-updates.json
  
  # Update policy label only via stdin
  echo '{"label":"updated-policy-name"}' | metalcloud-cli fw-policy update 456 --config-source pipe
  
  # Update policy action and rules
  metalcloud-cli firmware-policy update 789 --config-source new-config.json

```
metalcloud-cli firmware-policy update policy_id [flags]
```

### Options

```
      --config-source string   Source of the firmware policy configuration updates. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli firmware-policy](metalcloud-cli_firmware-policy.md)	 - Manage server firmware upgrade policies and global firmware configurations

