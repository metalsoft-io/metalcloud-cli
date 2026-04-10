## metalcloud-cli site dhcp-oob-reservations add

Add a DHCP OOB reservation entry

### Synopsis

Add a new MAC address to IP address mapping to the site's DHCP OOB reservations.

The IP address must belong to an OOB subnet assigned to the site. The MAC address
must be in standard format (e.g., AA:BB:CC:DD:EE:FF or aa:bb:cc:dd:ee:ff).

Entries can be added one at a time using --mac and --ip flags, or in bulk
using --from-json with a JSON string or --from-file with a path to a JSON file.
The JSON format is an object with MAC addresses as keys and IP addresses as values.

Examples:
  # Add a single reservation
  metalcloud-cli site dhcp-oob-reservations add site-01 --mac AA:BB:CC:DD:EE:FF --ip 10.0.0.100

  # Add multiple reservations from JSON string
  metalcloud-cli site dhcp-oob-reservations add site-01 --from-json '{"AA:BB:CC:DD:EE:FF":"10.0.0.100","11:22:33:44:55:66":"10.0.0.101"}'

  # Add multiple reservations from a JSON file
  metalcloud-cli site dhcp-oob-reservations add site-01 --from-file reservations.json

```
metalcloud-cli site dhcp-oob-reservations add site_id_or_name [flags]
```

### Options

```
      --from-file string   Path to a JSON file with MAC-to-IP mappings
      --from-json string   JSON object with MAC-to-IP mappings
  -h, --help               help for add
      --ip string          IP address to map to the MAC address
      --mac string         MAC address (e.g., AA:BB:CC:DD:EE:FF)
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

