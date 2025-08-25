## metalcloud-cli file-share update-hosts

Update hosts configuration for a file share

### Synopsis

Update the hosts configuration for an existing file share.

This command allows you to modify which hosts have access to the file share
and their mount configurations. You can add new hosts, remove existing ones,
or update their access permissions and mount settings.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to update hosts for

Required Flags:
  --config-source               Source of the file share hosts configuration
                               Can be 'pipe' for stdin input or path to a JSON file

Hosts Configuration Format:
The configuration should be a JSON object containing the hosts configuration:
- hosts: Array of host objects with their access settings
- mount_path: Default mount path for hosts
- access_permissions: Default access permissions (read, write, read-write)

Examples:
  # Update hosts configuration from a JSON file
  metalcloud-cli file-share update-hosts my-infrastructure 12345 --config-source hosts.json

  # Update using pipe input
  echo '{"hosts":[{"id":"host1","access":"read-write"}]}' | metalcloud-cli file-share update-hosts my-infrastructure 12345 --config-source pipe

  # Update with infrastructure ID
  metalcloud-cli file-share update-hosts 100 12345 --config-source /path/to/hosts.json

```
metalcloud-cli file-share update-hosts infrastructure_id_or_label file_share_id [flags]
```

### Options

```
      --config-source string   Source of the file share hosts configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-hosts
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

* [metalcloud-cli file-share](metalcloud-cli_file-share.md)	 - Manage file shares for infrastructure resources

