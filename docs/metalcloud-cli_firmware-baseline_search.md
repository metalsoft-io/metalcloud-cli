## metalcloud-cli firmware-baseline search

Search firmware baselines by criteria

### Synopsis

Search firmware baselines by criteria.

This command allows you to search for firmware baselines using specific criteria
such as vendor, datacenter, server type, and component filters.
Search criteria are provided through a configuration file (JSON or YAML format).

The search can filter baselines by:
- Vendor (e.g., DELL)
- Datacenter locations
- Server types and OS templates
- Component filters for specific hardware

Use the 'search-example' command to see available search criteria and their format.

Required Flags:
  --search-source    Source of the search criteria (JSON/YAML file path or 'pipe')

Examples:
  metalcloud-cli firmware-baseline search --search-source ./search-criteria.json
  cat search.json | metalcloud-cli fw-baseline search --search-source pipe
  metalcloud-cli baseline search --search-source ./dell-baselines.yaml

Search criteria example (search-criteria.json):
{
  "vendor": "DELL",
  "baselineFilter": {
    "datacenter": ["datacenter-1"],
    "serverType": ["dell_r740", "dell_r640"],
    "osTemplate": ["ubuntu-20.04"],
    "baselineId": ["baseline-1"]
  },
  "serverComponentFilter": {
    "dellComponentFilter": {
      "componentId": "component-1"
    }
  }
}

```
metalcloud-cli firmware-baseline search [flags]
```

### Options

```
  -h, --help                   help for search
      --search-source string   Source of the search criteria. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli firmware-baseline](metalcloud-cli_firmware-baseline.md)	 - Manage firmware baselines for consistent hardware configurations

