## metalcloud-cli firmware-baseline search-example

Display search criteria template for firmware baseline search

### Synopsis

Display search criteria template for firmware baseline search.

This command outputs a comprehensive example search criteria file that shows all available
search options for finding firmware baselines. The example includes all searchable fields
with their descriptions and sample values.

The search criteria template covers:
- Vendor filtering (DELL, etc.)
- Baseline filtering (datacenter, server type, OS template, baseline ID)
- Component filtering for specific hardware components

Use this template as a starting point for creating your own search criteria.
Copy the output to a file, modify the values as needed, and use it with the search command.

Examples:
  metalcloud-cli firmware-baseline search-example > search-template.json
  metalcloud-cli fw-baseline search-example | grep -A 10 "vendor"
  metalcloud-cli baseline search-example | jq '.baselineFilter'

```
metalcloud-cli firmware-baseline search-example [flags]
```

### Options

```
  -h, --help   help for search-example
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

