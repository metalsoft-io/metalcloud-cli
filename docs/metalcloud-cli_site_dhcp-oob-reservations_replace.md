## metalcloud-cli site dhcp-oob-reservations replace

Replace all DHCP OOB reservation entries

### Synopsis

Replace the entire DHCP Option82 to IP address mapping with the provided entries.

This removes all existing reservations and sets the mapping to the provided JSON
object. All IP addresses must belong to OOB subnets assigned to the site.
Input can be a JSON string via --from-json or a JSON file via --from-file.

Examples:
  # Replace all reservations from JSON string
  metalcloud-cli site dhcp-oob-reservations replace site-01 --from-json '{"AA:BB:CC:DD:EE:FF":"10.0.0.100","11:22:33:44:55:66":"10.0.0.101"}'

  # Replace all reservations from a file
  metalcloud-cli site dhcp-oob-reservations replace site-01 --from-file reservations.json

  # Clear all reservations
  metalcloud-cli site dhcp-oob-reservations replace site-01 --from-json '{}'

```
metalcloud-cli site dhcp-oob-reservations replace site_id_or_name [flags]
```

### Options

```
      --from-file string   Path to a JSON file with MAC-to-IP mappings
      --from-json string   JSON object with MAC-to-IP mappings
  -h, --help               help for replace
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

