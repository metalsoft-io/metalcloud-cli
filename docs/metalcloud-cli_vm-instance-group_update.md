## metalcloud-cli vm-instance-group update

Update VM instance group configuration

### Synopsis

Update VM instance group configuration.

This command allows you to modify the configuration of an existing VM instance
group. You can update the label or custom variables associated with the group.
The current configuration revision is fetched automatically and sent as the
If-Match header, so concurrent modifications are rejected by the API.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Optional Flags:
  --label string                    Set or update the VM instance group label
  --custom-variables-source string  Source of the custom variables.
                                    Can be 'pipe' or path to a JSON file.

Flag Dependencies:
  At least one of --label or --custom-variables-source must be provided

Examples:
  # Update the label of a VM instance group
  metalcloud-cli vm-instance-group update 12345 67890 --label "Web Servers"

  # Update custom variables from a JSON file
  metalcloud-cli vmg edit my-infra 67890 --custom-variables-source /path/to/vars.json

  # Update custom variables from stdin
  echo '{"env": "production"}' | metalcloud-cli vm-group update 12345 67890 --custom-variables-source pipe

```
metalcloud-cli vm-instance-group update infrastructure_id_or_label vm_instance_group_id [flags]
```

### Options

```
      --custom-variables-source string   Source of the custom variables. Can be 'pipe' or path to a JSON file.
  -h, --help                             help for update
      --label string                     Set the VM instance group label.
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

* [metalcloud-cli vm-instance-group](metalcloud-cli_vm-instance-group.md)	 - Manage VM instance groups within infrastructures

