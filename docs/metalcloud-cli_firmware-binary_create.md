## metalcloud-cli firmware-binary create

Create a new firmware binary from configuration file

### Synopsis

Create a new firmware binary from configuration file.

This command creates a new firmware binary registration in the system using a configuration
file that defines all the binary's properties, compatibility information, and metadata.

The firmware binary will be validated for:
- Hardware compatibility specifications
- Version information consistency
- Required metadata completeness
- Update severity classification

Use the 'config-example' command to generate a template configuration file with all
available options and their descriptions.

Required Flags:
  --config-source    Source of the firmware binary configuration (JSON/YAML file path or 'pipe')

The configuration file must include:
- Basic metadata (name, catalogId)
- Vendor download URL
- Hardware compatibility information (vendorSupportedDevices, vendorSupportedSystems)
- Update requirements (rebootRequired, updateSeverity)

Optional configuration includes:
- Package identification (packageId, packageVersion)
- External references (externalId, vendorInfoUrl, cacheDownloadUrl)
- Release information (vendorReleaseTimestamp)
- Vendor details

Examples:
  metalcloud-cli firmware-binary create --config-source ./bios-update.json
  cat firmware-config.json | metalcloud-cli fw-binary create --config-source pipe
  metalcloud-cli firmware-bin new --config-source ./dell-r640-bios.yaml

Configuration file example (bios-update.json):
{
  "name": "BIOS-R740-2.15.0",
  "catalogId": 1,
  "vendorDownloadUrl": "https://dell.com/downloads/firmware/R740/BIOS-2.15.0.bin",
  "vendorInfoUrl": "https://dell.com/support/firmware/R740/BIOS/2.15.0",
  "externalId": "DELL-R740-BIOS-2.15.0",
  "packageId": "BIOS",
  "packageVersion": "2.15.0",
  "rebootRequired": true,
  "updateSeverity": "recommended",
  "vendorReleaseTimestamp": "2024-04-01T12:00:00Z",
  "vendorSupportedDevices": [
    {
      "model": "PowerEdge R740",
      "type": "server"
    }
  ],
  "vendorSupportedSystems": [
    {
      "os": "any",
      "version": "any"
    }
  ],
  "vendor": {
    "name": "Dell Inc.",
    "contact": "support@dell.com"
  }
}

```
metalcloud-cli firmware-binary create [flags]
```

### Options

```
      --config-source string   Source of the new firmware binary configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli firmware-binary](metalcloud-cli_firmware-binary.md)	 - Manage individual firmware binary files and packages

