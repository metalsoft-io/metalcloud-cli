## metalcloud-cli firmware-binary config-example

Display configuration file template for creating firmware binaries

### Synopsis

Display configuration file template for creating firmware binaries.

This command outputs a comprehensive example configuration file that shows all available
options for creating firmware binaries. The example includes all required and optional
fields with their descriptions and sample values.

The configuration template covers:
- Basic metadata (name, catalogId, packageId, packageVersion)
- File information (vendorDownloadUrl, vendorInfoUrl, cacheDownloadUrl)
- Hardware compatibility (vendorSupportedDevices, vendorSupportedSystems)
- Installation requirements (rebootRequired, updateSeverity)
- Vendor information and external references

Use this template as a starting point for creating your own firmware binary configurations.
Copy the output to a file, modify the values as needed, and use it with the create command.

Examples:
  metalcloud-cli firmware-binary config-example > my-firmware.json
  metalcloud-cli fw-binary config-example | grep -A 50 "dell"

```
metalcloud-cli firmware-binary config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

