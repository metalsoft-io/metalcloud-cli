## metalcloud-cli network-configuration device-template import-library

Bulk-import a directory of templates as a single library

### Synopsis

Bulk-import every network device configuration template descriptor found in a
directory, grouping them all under a single library label.

Each file in the directory is one template descriptor (JSON or YAML) with the
same fields as 'config-example' - the preparation and configuration fields are
base64-encoded commands. Files with a .json, .yaml or .yml extension are imported
in name order; any file's own libraryLabel is overridden with <library_label> so
the whole directory forms one library. A file that cannot be read or parsed is
reported and skipped so one bad file does not abort the rest.

Arguments:
  library_label   The library label to assign to every imported template

Required Flags:
  --dir           Directory holding the template descriptor files

Examples:
  # Preview what would be imported
  metalcloud-cli network-configuration device-template import-library spectrumx --dir ./templates --dry-run

  # Import every descriptor in ./templates as the 'spectrumx' library
  metalcloud-cli nc dt import-library spectrumx --dir ./templates

```
metalcloud-cli network-configuration device-template import-library <library_label> [flags]
```

### Options

```
      --dir string   Directory holding the template descriptor files (*.json, *.yaml, *.yml).
      --dry-run      Report what would be imported without creating anything.
  -h, --help         help for import-library
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

* [metalcloud-cli network-configuration device-template](metalcloud-cli_network-configuration_device-template.md)	 - Manage network devices configuration templates

