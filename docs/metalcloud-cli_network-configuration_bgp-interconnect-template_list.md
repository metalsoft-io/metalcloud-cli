## metalcloud-cli network-configuration bgp-interconnect-template list

List network device BGP interconnect configuration templates with optional filtering

### Synopsis

List all network device BGP interconnect configuration templates with optional filtering.

This command displays all network device BGP interconnect configuration templates that are registered in the system.
You can filter the results by network device driver to focus on specific groups of templates.
Flags:
  --filter-network-device-driver   Filter templates by network device driver

Examples:
  # List all network device BGP interconnect configuration templates (default)
  metalcloud-cli network-configuration bgp-interconnect-template list

  # List templates for a specific network device driver
  metalcloud-cli network-configuration bgp-interconnect-template list --filter-network-device-driver junos

```
metalcloud-cli network-configuration bgp-interconnect-template list [flags]
```

### Options

```
      --filter-id strings                      Filter by template ID.
      --filter-network-device-driver strings   Filter by network device driver.
  -h, --help                                   help for list
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

* [metalcloud-cli network-configuration bgp-interconnect-template](metalcloud-cli_network-configuration_bgp-interconnect-template.md)	 - Manage network devices BGP interconnect configuration templates

