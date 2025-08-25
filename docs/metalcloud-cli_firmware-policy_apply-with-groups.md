## metalcloud-cli firmware-policy apply-with-groups

Apply all firmware policies linked to server instance groups

### Synopsis

Apply all firmware policies that are linked to server instance groups.

This command executes all active firmware policies that have server instance groups
associated with them. It will only affect servers that belong to the specified
server instance groups in each policy's configuration.

The command respects the global firmware configuration settings for timing and
scheduling constraints. If global firmware upgrades are disabled or outside
the configured time window, the command may be blocked.

No flags or arguments are required for this command.

Prerequisites:
  - At least one firmware policy must exist with server instance groups assigned
  - Global firmware configuration must allow policy execution
  - Servers in the target groups must be accessible and eligible for firmware updates

Examples:
  # Apply all policies linked to server instance groups
  metalcloud-cli firmware-policy apply-with-groups
  
  # Apply policies using alias
  metalcloud-cli fw-policy apply-with-groups
  
  # Check global config before applying
  metalcloud-cli firmware-policy global-config get
  metalcloud-cli firmware-policy apply-with-groups

```
metalcloud-cli firmware-policy apply-with-groups [flags]
```

### Options

```
  -h, --help   help for apply-with-groups
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

