## metalcloud-cli firmware-baseline update

Update an existing firmware baseline

### Synopsis

Update an existing firmware baseline.

This command allows you to update the configuration of an existing firmware baseline.
Updates are provided through a configuration file (JSON or YAML format) that contains
the new settings to apply.

The configuration file should contain only the fields you want to update. Common
updates include:
- Basic metadata (name, description)
- Level and level filter specifications  
- Catalog associations

The baseline will be revalidated after updates to ensure consistency and compatibility.

Required Flags:
  --config-source    Source of the configuration updates (JSON/YAML file path or 'pipe')

Arguments:
  firmware_baseline_id    The ID of the firmware baseline to update

Examples:
  metalcloud-cli firmware-baseline update 54321 --config-source ./baseline-updates.json
  cat updates.json | metalcloud-cli fw-baseline update production-baseline --config-source pipe
  metalcloud-cli baseline edit dell-r640-standard --config-source ./version-update.yaml

Configuration file example (baseline-updates.json):
{
  "name": "Updated Production Dell R640 Baseline",
  "description": "Updated with latest security patches",
  "level": "PRODUCTION",
  "levelFilter": ["dell_r640", "dell_r640_gen2"],
  "catalog": ["dell-catalog-r640-v2", "dell-catalog-common"]
}

```
metalcloud-cli firmware-baseline update firmware_baseline_id [flags]
```

### Options

```
      --config-source string   Source of the firmware baseline configuration updates. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

