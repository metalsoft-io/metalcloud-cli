## metalcloud-cli firmware-baseline

Manage firmware baselines for consistent hardware configurations

### Synopsis

Manage firmware baselines for consistent hardware configurations.

Firmware baselines define standardized firmware configurations for specific hardware
types and deployment scenarios. They specify the firmware level and filtering criteria
for consistent hardware management across your infrastructure.

A firmware baseline includes:
  • Name and description for identification
  • Level specification (e.g., PRODUCTION, DEVELOPMENT)
  • Level filter for targeting specific hardware types
  • Catalog associations for firmware sources

Use cases:
  • Standardizing firmware configurations across server fleets
  • Defining deployment levels for different environments
  • Managing firmware catalog associations
  • Creating hardware-specific configuration templates

### Options

```
  -h, --help   help for firmware-baseline
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
* [metalcloud-cli firmware-baseline config-example](metalcloud-cli_firmware-baseline_config-example.md)	 - Display configuration file template for creating firmware baselines
* [metalcloud-cli firmware-baseline create](metalcloud-cli_firmware-baseline_create.md)	 - Create a new firmware baseline from configuration file
* [metalcloud-cli firmware-baseline delete](metalcloud-cli_firmware-baseline_delete.md)	 - Delete a firmware baseline permanently
* [metalcloud-cli firmware-baseline get](metalcloud-cli_firmware-baseline_get.md)	 - Get detailed information about a specific firmware baseline
* [metalcloud-cli firmware-baseline list](metalcloud-cli_firmware-baseline_list.md)	 - List all firmware baselines
* [metalcloud-cli firmware-baseline search](metalcloud-cli_firmware-baseline_search.md)	 - Search firmware baselines by criteria
* [metalcloud-cli firmware-baseline search-example](metalcloud-cli_firmware-baseline_search-example.md)	 - Display search criteria template for firmware baseline search
* [metalcloud-cli firmware-baseline update](metalcloud-cli_firmware-baseline_update.md)	 - Update an existing firmware baseline

