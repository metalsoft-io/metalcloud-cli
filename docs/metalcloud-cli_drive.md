## metalcloud-cli drive

Manage drives within infrastructures

### Synopsis

Manage drives within infrastructures including creation, configuration, metadata updates, and host assignments.

Drives are storage devices that can be attached to servers within an infrastructure. This command group
provides comprehensive drive management capabilities including listing, creating, updating configurations,
managing metadata, and controlling host assignments.

Available Commands:
  list          List all drives for an infrastructure
  get           Get detailed information about a specific drive
  create        Create a new drive with specified configuration
  delete        Remove a drive from the infrastructure
  update-config Update drive configuration settings
  update-meta   Update drive metadata
  get-hosts     Show hosts assigned to a drive
  update-hosts  Update host assignments for a drive
  config-info   Get configuration information for a drive

Examples:
  # List all drives in an infrastructure
  metalcloud-cli drive list my-infrastructure

  # Get details of a specific drive
  metalcloud-cli drive get my-infrastructure 12345

  # Create a new drive from JSON configuration
  metalcloud-cli drive create my-infrastructure --config-source drive-config.json

### Options

```
  -h, --help   help for drive
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
* [metalcloud-cli drive config-info](metalcloud-cli_drive_config-info.md)	 - Get configuration information for a drive
* [metalcloud-cli drive create](metalcloud-cli_drive_create.md)	 - Create a new drive with specified configuration
* [metalcloud-cli drive delete](metalcloud-cli_drive_delete.md)	 - Remove a drive from the infrastructure
* [metalcloud-cli drive get](metalcloud-cli_drive_get.md)	 - Get detailed information about a specific drive
* [metalcloud-cli drive get-hosts](metalcloud-cli_drive_get-hosts.md)	 - Show hosts assigned to a drive
* [metalcloud-cli drive list](metalcloud-cli_drive_list.md)	 - List all drives within an infrastructure
* [metalcloud-cli drive update-config](metalcloud-cli_drive_update-config.md)	 - Update drive configuration settings
* [metalcloud-cli drive update-hosts](metalcloud-cli_drive_update-hosts.md)	 - Update host assignments for a drive
* [metalcloud-cli drive update-meta](metalcloud-cli_drive_update-meta.md)	 - Update drive metadata

