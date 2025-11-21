## metalcloud-cli network-device-configuration-template list

List network device configuration templates with optional filtering

### Synopsis

List all network device configuration templates with optional filtering.

This command displays all network device configuration templates that are registered in the system.
You can filter the results by library label to focus on specific groups of templates.
Flags:
  --filter-library-label   Filter templates by library label

Examples:
  # List all network device configuration templates (default)
  metalcloud-cli network-device-configuration-template list

  # List templates with a specific library label
  metalcloud-cli network-device-configuration-template list --filter-library-label example-label

```
metalcloud-cli network-device-configuration-template list [flags]
```

### Options

```
      --filter-id strings              Filter by template ID.
      --filter-library-label strings   Filter by template library label.
  -h, --help                           help for list
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

* [metalcloud-cli network-device-configuration-template](metalcloud-cli_network-device-configuration-template.md)	 - Manage network devices configuration templates

