## metalcloud-cli server firmware components

List firmware components for a server

### Synopsis

List firmware components for a server.

This command displays all firmware components available on the specified server,
including their IDs, names, class, vendor, model, versions and, for the components of a
GPU baseboard (an HGX tray), the baseboard they belong to.

Required Arguments:
  server_id              The ID of the server to query

Optional Flags:
  --class                Only the components of these classes (bios, bmc, gpu, nic, hgx_hmc, ...)
  --baseboard            Only the components of these GPU baseboards

Examples:
  # List firmware components for server with ID 123
  metalcloud-cli server firmware components 123

  # Only the GPUs
  metalcloud-cli server firmware components 123 --class gpu

  # Using short alias
  metalcloud-cli server fw ls 123


```
metalcloud-cli server firmware components server_id [flags]
```

### Options

```
      --baseboard strings   Only the components of these GPU baseboards.
      --class strings       Only the components of these classes.
  -h, --help                help for components
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

