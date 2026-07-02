## metalcloud-cli network-configuration device-template export-library

Export all templates of a single library to a directory

### Synopsis

Export every network device configuration template that belongs to a library to
a directory, one descriptor file per template.

Each file is the inverse of 'import-library' input - the create fields only, no
id or timestamps - so an exported directory can be re-imported as-is (optionally
under a different library label). The output directory is created if it does not
exist; files are named 'template-<id>.json'.

Arguments:
  library_label   The library label whose templates should be exported

Required Flags:
  --dir           Directory to write the descriptor files into

Examples:
  # Export the 'spectrumx' library to ./spectrumx-export
  metalcloud-cli network-configuration device-template export-library spectrumx --dir ./spectrumx-export

  # Round-trip: export, then re-import under a new label
  metalcloud-cli nc dt export-library spectrumx --dir ./lib
  metalcloud-cli nc dt import-library spectrumx-copy --dir ./lib

```
metalcloud-cli network-configuration device-template export-library <library_label> [flags]
```

### Options

```
      --dir string   Directory to write the exported descriptor files into (created if missing).
  -h, --help         help for export-library
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

