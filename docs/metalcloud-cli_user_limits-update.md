## metalcloud-cli user limits-update

Update resource limits for a specific user

### Synopsis

Update the resource limits for a specific user account to control their resource allocation.

This command allows you to modify compute node, drive, and infrastructure limits that restrict
how many resources the user can provision. Changes take effect immediately.

Arguments:
  user_id                 The numeric ID of the user whose limits to update

Required Flags:
  --config-source         Source of user limits configuration (JSON file path or 'pipe')

Configuration File Format (JSON):
  {
    "computeNodesInstancesToProvisionLimit": 100,
    "drivesAttachedToInstancesLimit": 200,
    "infrastructuresLimit": 10
  }

Examples:
  metalcloud-cli user limits-update 12345 --config-source limits.json
  echo '{"computeNodesInstancesToProvisionLimit": 50}' | metalcloud-cli user limits-update 12345 --config-source pipe

```
metalcloud-cli user limits-update user_id [flags]
```

### Options

```
      --config-source string   Source of the user limits configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for limits-update
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

