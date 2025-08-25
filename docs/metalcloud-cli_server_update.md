## metalcloud-cli server update

Update server information

### Synopsis

Update server information.

This command updates server configuration using a JSON configuration file or 
piped JSON data. The configuration must be provided via the --config-source flag.

Required Arguments:
  server_id              The ID of the server to update

Required Flags:
  --config-source        Source of the server update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update server using JSON configuration file
  metalcloud-cli server update 123 --config-source ./server-update.json

  # Update server using piped JSON configuration
  echo '{"vendor": "Dell", "model": "PowerEdge R740"}' | metalcloud-cli server update 123 --config-source pipe


```
metalcloud-cli server update server_id [flags]
```

### Options

```
      --config-source string   Source of the server update configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

