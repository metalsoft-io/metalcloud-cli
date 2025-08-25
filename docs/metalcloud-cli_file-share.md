## metalcloud-cli file-share

Manage file shares for infrastructure resources

### Synopsis

Manage file shares for infrastructure resources including creating, updating, deleting, and configuring shared storage.

File shares provide shared storage capabilities across multiple instances within an infrastructure.
Use these commands to manage the lifecycle and configuration of file shares.

Available Commands:
  list           List all file shares for an infrastructure
  get            Get detailed information about a specific file share
  create         Create a new file share with specified configuration
  delete         Delete an existing file share
  update-config  Update file share configuration
  update-meta    Update file share metadata
  get-hosts      Get hosts configured for a file share
  update-hosts   Update hosts configuration for a file share
  config-info    Get configuration information for a file share

Examples:
  # List all file shares for an infrastructure
  metalcloud-cli file-share list my-infrastructure

  # Get details of a specific file share
  metalcloud-cli file-share get my-infrastructure 12345

  # Create a new file share from a configuration file
  metalcloud-cli file-share create my-infrastructure --config-source config.json

### Options

```
  -h, --help   help for file-share
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli file-share config-info](metalcloud-cli_file-share_config-info.md)	 - Get configuration information for a file share
* [metalcloud-cli file-share create](metalcloud-cli_file-share_create.md)	 - Create a new file share with specified configuration
* [metalcloud-cli file-share delete](metalcloud-cli_file-share_delete.md)	 - Delete an existing file share
* [metalcloud-cli file-share get](metalcloud-cli_file-share_get.md)	 - Get detailed information about a specific file share
* [metalcloud-cli file-share get-hosts](metalcloud-cli_file-share_get-hosts.md)	 - Get hosts configured for a file share
* [metalcloud-cli file-share list](metalcloud-cli_file-share_list.md)	 - List all file shares for an infrastructure
* [metalcloud-cli file-share update-config](metalcloud-cli_file-share_update-config.md)	 - Update file share configuration
* [metalcloud-cli file-share update-hosts](metalcloud-cli_file-share_update-hosts.md)	 - Update hosts configuration for a file share
* [metalcloud-cli file-share update-meta](metalcloud-cli_file-share_update-meta.md)	 - Update file share metadata

