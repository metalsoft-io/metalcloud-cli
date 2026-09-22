## metalcloud-cli network-device snmp-service enable

Enable the SNMP agent on a network device

### Synopsis

Enable the SNMP agent running on the network device, optionally overriding the
listening port, the community string and the contact information.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Optional Flags:
  --port        SNMP listening port
  --community   SNMP community string
  --contact     SNMP contact information

Examples:
  # Enable the SNMP agent with the stored defaults
  metalcloud-cli network-device snmp-service enable 12345

  # Enable it on a custom port and community
  metalcloud-cli network-device snmp-service enable 12345 --port 1161 --community public

```
metalcloud-cli network-device snmp-service enable <network_device_id> [flags]
```

### Options

```
      --community string   SNMP community string.
      --contact string     SNMP contact information.
  -h, --help               help for enable
      --port int32         SNMP listening port. (default 161)
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

* [metalcloud-cli network-device snmp-service](metalcloud-cli_network-device_snmp-service.md)	 - Manage the SNMP agent running on a network device

