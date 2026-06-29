## metalcloud-cli fabric rescan-links

Re-scan (discover) the fabric's links from LLDP

### Synopsis

Re-scan the links of a network fabric, deriving the physical-link records from
the LLDP data that flows once the ports are up. This is the "Discover links"
action: run it after the first fabric deploy, then deploy again so the links
become MetalSoft-managed.

It is idempotent - on a fabric whose links already exist it matches rows in
place without creating duplicates.

Arguments:
  fabric_id    The ID or label of the fabric whose links to rescan

Optional Flags:
  --update-lldp   Update LLDP information during the rescan (default true).

Examples:
  metalcloud-cli fabric rescan-links 12345
  metalcloud-cli fabric discover-links my-fabric
  metalcloud-cli fabric rescan-links 12345 --update-lldp=false

```
metalcloud-cli fabric rescan-links fabric_id [flags]
```

### Options

```
  -h, --help          help for rescan-links
      --update-lldp   Update LLDP information during the link rescan. (default true)
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

