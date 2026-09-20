## metalcloud-cli network-device snmp-service

Manage the SNMP agent running on a network device

### Synopsis

Enable or disable the SNMP agent running on the network device itself.

Both operations are asynchronous and return the job that carries them out.

### Options

```
  -h, --help   help for snmp-service
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
* [metalcloud-cli network-device snmp-service disable](metalcloud-cli_network-device_snmp-service_disable.md)	 - Disable the SNMP agent on a network device
* [metalcloud-cli network-device snmp-service enable](metalcloud-cli_network-device_snmp-service_enable.md)	 - Enable the SNMP agent on a network device

