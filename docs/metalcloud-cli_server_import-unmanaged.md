## metalcloud-cli server import-unmanaged

Import a server whose lifecycle MetalSoft does not manage

### Synopsis

Import a server that MetalSoft does not manage the lifecycle of, using an
externally supplied hardware description.

The configuration must describe the site, the server type and the server
interfaces together with the network device ports they are cabled to.

Required Flags:
  --config-source        Source of the import configuration. Can be 'pipe' or path to a JSON file.

Examples:
  # Import an unmanaged server from a JSON file
  metalcloud-cli server import-unmanaged --config-source ./unmanaged-server.json

  # Import an unmanaged server from piped configuration
  metalcloud-cli server config-example import-unmanaged | metalcloud-cli server import-unmanaged --config-source pipe


```
metalcloud-cli server import-unmanaged [flags]
```

### Options

```
      --config-source string   Source of the import configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for import-unmanaged
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

