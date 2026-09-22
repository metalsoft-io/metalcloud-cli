## metalcloud-cli dhcp-reservation create

Create a DHCP reservation

### Synopsis

Create a DHCP reservation in a site for one IP version.

The reservation can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source.

A manual allocation pins one IP: pass --ip together with --subnet-id. An auto
allocation takes the next free IP of one or more OOB subnet pools: pass
--subnet-pool-id (repeatable). The two forms are mutually exclusive.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'

Required Flags (one of):
  --ip + --subnet-id   Pin this IP from this subnet (manual allocation)
  --subnet-pool-id     Allocate from these pools, in order (auto allocation)
  --config-source      'pipe' to read from stdin, or a path to a JSON/YAML file

Optional Flags:
  --mac-address   MAC address to match, as six colon-separated hex octets
  --circuit-id    Option 82 circuit id to match
  --device-type   Device type guard/filter
  --vendor        Vendor guard/filter
  --giaddr        Relay gateway address (giaddr) to match
  --priority      Tiebreaker among equally specific matches (lower wins)

Examples:
  metalcloud-cli dhcp-reservation create 1 ipv4 --mac-address AA:BB:CC:DD:EE:FF --ip 192.168.1.10 --subnet-id 3
  metalcloud-cli dhcp-reservation create 1 ipv4 --circuit-id leaf-01:swp1 --subnet-pool-id 3 --subnet-pool-id 4
  metalcloud-cli dhcp-reservation create 1 ipv4 --config-source reservation.yaml

```
metalcloud-cli dhcp-reservation create site_id_or_label ip_version [flags]
```

### Options

```
      --circuit-id string      Option 82 circuit id to match.
      --config-source string   Source of the new DHCP reservation configuration. Can be 'pipe' or path to a JSON/YAML file.
      --device-type string     Device type guard/filter.
      --giaddr string          Relay gateway address (giaddr) to match.
  -h, --help                   help for create
      --ip string              Manual allocation: the IP address to pin for this reservation.
      --mac-address string     MAC address to match, as six colon-separated hex octets.
      --priority int           Tiebreaker among equally specific matches (lower = higher precedence).
      --subnet-id int          Manual allocation: the OOB subnet id that contains the pinned IP.
      --subnet-pool-id ints    Auto allocation: OOB subnet pool ids to allocate the next free IP from, in order.
      --vendor string          Vendor guard/filter.
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

* [metalcloud-cli dhcp-reservation](metalcloud-cli_dhcp-reservation.md)	 - Manage site DHCP reservations

