## metalcloud-cli network-device vendor

Manage the network device vendor profiles

### Synopsis

Manage the per-vendor profiles that drive SNMP monitoring, health checks and
configuration backups for each network device driver.

### Options

```
  -h, --help   help for vendor
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

* [metalcloud-cli network-device](metalcloud-cli_network-device.md)	 - Manage network devices (switches) in the infrastructure
* [metalcloud-cli network-device vendor config-example](metalcloud-cli_network-device_vendor_config-example.md)	 - Example configuration for updating a vendor profile
* [metalcloud-cli network-device vendor get](metalcloud-cli_network-device_vendor_get.md)	 - Get one network device vendor profile
* [metalcloud-cli network-device vendor list](metalcloud-cli_network-device_vendor_list.md)	 - List the network device vendor profiles
* [metalcloud-cli network-device vendor update](metalcloud-cli_network-device_vendor_update.md)	 - Update one network device vendor profile

