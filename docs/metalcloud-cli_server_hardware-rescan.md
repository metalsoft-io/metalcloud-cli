## metalcloud-cli server hardware-rescan

Re-read the hardware inventory of a server

### Synopsis

Re-read the hardware inventory of a server.

The rescan updates the recorded hardware configuration of the server. By default
the server is not rebooted, which preserves the existing switch connections; pass
--reboot-allowed to let MetalSoft reboot the server and run a full LLDP interface
discovery.

Required Arguments:
  server_id              The ID of the server to rescan

Optional Flags:
  --reboot-allowed       Allow the server to be rebooted for full LLDP interface discovery

Examples:
  # Rescan the hardware of server 123 without rebooting it
  metalcloud-cli server hardware-rescan 123

  # Rescan the hardware of server 123, allowing a reboot
  metalcloud-cli server hardware-rescan 123 --reboot-allowed


```
metalcloud-cli server hardware-rescan server_id [flags]
```

### Options

```
  -h, --help             help for hardware-rescan
      --reboot-allowed   Allow the server to be rebooted for full LLDP interface discovery.
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

