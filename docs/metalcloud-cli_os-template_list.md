## metalcloud-cli os-template list

List all available OS templates

### Synopsis

List all available OS templates in the system.

This command displays a table of all OS templates with their basic information
including ID, name, label, device type, status, visibility, and timestamps.

The output includes:
  - Template ID (unique identifier)
  - Name (human-readable template name)
  - Label (unique template label)
  - Device Type (server, switch, etc.)
  - Status (ready, active, used, archived)
  - Visibility (public, private)
  - Created/Modified timestamps

Examples:
  # List all OS templates
  metalcloud-cli os-template list
  
  # List templates using alias
  metalcloud-cli templates ls

```
metalcloud-cli os-template list [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli os-template](metalcloud-cli_os-template.md)	 - Manage OS templates for server deployments

