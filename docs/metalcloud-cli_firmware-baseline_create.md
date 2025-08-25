## metalcloud-cli firmware-baseline create

Create a new firmware baseline from configuration file

### Synopsis

Create a new firmware baseline from configuration file.

This command creates a new firmware baseline definition using a configuration file that
specifies all the baseline properties and target hardware specifications.

The firmware baseline will be validated for:
- Required field completeness
- Level and level filter consistency
- Catalog reference validity

Use the 'config-example' command to generate a template configuration file with all
available options and their descriptions.

Required Flags:
  --config-source    Source of the firmware baseline configuration (JSON/YAML file path or 'pipe')

The configuration file must include:
- Basic metadata (name, level, levelFilter)

Optional configuration includes:
- Description and catalog associations

Examples:
  metalcloud-cli firmware-baseline create --config-source ./production-baseline.json
  cat baseline-config.json | metalcloud-cli fw-baseline create --config-source pipe
  metalcloud-cli baseline new --config-source ./dell-r640-baseline.yaml

Configuration file example (production-baseline.json):
{
  "name": "Production Dell R640 Baseline",
  "description": "Standard firmware configuration for Dell PowerEdge R640 production servers",
  "level": "PRODUCTION",
  "levelFilter": ["dell_r640", "dell_r640_gen2"],
  "catalog": ["dell-catalog-r640", "dell-catalog-common"]
}

```
metalcloud-cli firmware-baseline create [flags]
```

### Options

```
      --config-source string   Source of the new firmware baseline configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

