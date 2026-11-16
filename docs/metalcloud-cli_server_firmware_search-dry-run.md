## metalcloud-cli server firmware search-dry-run

Explain which firmware binaries match each component of a server

### Synopsis

Explain which firmware binaries match each component of a server.

Runs the binary search a firmware upgrade uses and shows, for each component, the binaries
that matched (with the rules that matched them) and the ones rejected (with the reason).
Nothing is stored or scheduled.

Required Arguments:
  server_id              The ID of the server

Optional Flags:
  --baseline-id          Search this firmware baseline
  --catalog-ids          Search these firmware catalogs (instead of a baseline)
  --all-versions         Include every matching version, not only the latest

Without --baseline-id or --catalog-ids the server's target baseline is searched.

Examples:
  # Search the server's target baseline
  metalcloud-cli server firmware search-dry-run 123

  # Search two catalogs, every version
  metalcloud-cli server firmware search-dry-run 123 --catalog-ids 10,11 --all-versions


```
metalcloud-cli server firmware search-dry-run server_id [flags]
```

### Options

```
      --all-versions             Include every matching version, not only the latest.
      --baseline-id int          Search this firmware baseline.
      --catalog-ids int64Slice   Search these firmware catalogs. (default [])
  -h, --help                     help for search-dry-run
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

* [metalcloud-cli server firmware](metalcloud-cli_server_firmware.md)	 - Server firmware management

