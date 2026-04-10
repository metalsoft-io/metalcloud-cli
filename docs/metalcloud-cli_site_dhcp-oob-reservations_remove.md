## metalcloud-cli site dhcp-oob-reservations remove

Remove a DHCP OOB reservation entry by MAC address

### Synopsis

Remove one or more MAC address entries from the site's DHCP OOB reservations.

MAC addresses can be specified with --mac flags, or loaded from a JSON file via
--from-file. The file should contain either a JSON array of MAC address strings
or a JSON object whose keys are MAC addresses.

Examples:
  # Remove a single reservation
  metalcloud-cli site dhcp-oob-reservations remove site-01 --mac AA:BB:CC:DD:EE:FF

  # Remove multiple reservations
  metalcloud-cli site dhcp-oob-reservations remove site-01 --mac AA:BB:CC:DD:EE:FF --mac 11:22:33:44:55:66

  # Remove reservations listed in a file (array format)
  metalcloud-cli site dhcp-oob-reservations remove site-01 --from-file remove-macs.json

  # Remove reservations listed in a file (object format - keys are used)
  metalcloud-cli site dhcp-oob-reservations remove site-01 --from-file reservations.json

```
metalcloud-cli site dhcp-oob-reservations remove site_id_or_name [flags]
```

### Options

```
      --from-file string   Path to a JSON file with MAC addresses (array or object keys)
  -h, --help               help for remove
      --mac strings        MAC address(es) to remove (can be specified multiple times)
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

* [metalcloud-cli site dhcp-oob-reservations](metalcloud-cli_site_dhcp-oob-reservations.md)	 - Manage DHCP Option82 OOB IP reservations for a site

