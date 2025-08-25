## metalcloud-cli server-cleanup-policy create

Create a new server cleanup policy

### Synopsis

Create a new server cleanup policy with specified configuration.

This command creates a new server cleanup policy that defines automated maintenance
procedures for servers. The policy configuration includes cleanup behaviors for
drives, RAID settings, and embedded NICs.

Required Flags:
  --label                                Label for the server cleanup policy
  --cleanup-drives                       Enable cleanup drives for OOB enabled servers
  --recreate-raid                        Enable RAID recreation
  --disable-embedded-nics                Enable disabling embedded NICs
  --raid-one-drive                       RAID configuration for single drive (e.g. "raid0")
  --raid-two-drives                      RAID configuration for two drives (e.g. "raid1")
  --raid-even-drives                     RAID configuration for even number of drives (>2)
  --raid-odd-drives                      RAID configuration for odd number of drives (>1)
  --skip-raid-actions                    Comma-separated list of RAID actions to skip

Required Permissions:
  - server_cleanup_policies:write

Examples:
  # Create a basic cleanup policy
  metalcloud-cli server-cleanup-policy create --label "basic-cleanup" \
    --cleanup-drives --recreate-raid --disable-embedded-nics \
    --raid-one-drive "raid0" --raid-two-drives "raid1" \
    --raid-even-drives "raid10" --raid-odd-drives "raid5" \
    --skip-raid-actions "cleanup"

```
metalcloud-cli server-cleanup-policy create [flags]
```

### Options

```
      --cleanup-drives             Enable cleanup drives for OOB enabled servers
      --disable-embedded-nics      Enable disabling embedded NICs
  -h, --help                       help for create
      --label string               Label for the server cleanup policy
      --raid-even-drives string    RAID configuration for even number of drives (>2)
      --raid-odd-drives string     RAID configuration for odd number of drives (>1)
      --raid-one-drive string      RAID configuration for single drive
      --raid-two-drives string     RAID configuration for two drives
      --recreate-raid              Enable RAID recreation
      --skip-raid-actions string   Comma-separated list of RAID actions to skip
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

* [metalcloud-cli server-cleanup-policy](metalcloud-cli_server-cleanup-policy.md)	 - Manage server cleanup policies for automated server maintenance

