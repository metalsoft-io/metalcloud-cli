## metalcloud-cli server-cleanup-policy update

Update an existing server cleanup policy

### Synopsis

Update an existing server cleanup policy by its ID.

This command updates the configuration of an existing server cleanup policy.
Only the flags that are provided will be updated, other settings remain unchanged.

Arguments:
  policy-id    The unique identifier of the server cleanup policy to update.
               This must be the numeric ID of the policy.

Optional Flags:
  --label                                New label for the server cleanup policy
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
  # Update policy label only
  metalcloud-cli server-cleanup-policy update 123 --label "updated-policy"

  # Update multiple settings
  metalcloud-cli scp update 123 --cleanup-drives --recreate-raid

  # Update RAID configurations
  metalcloud-cli srv-cp update 123 --raid-one-drive "RAID1" --raid-two-drives "RAID10"

```
metalcloud-cli server-cleanup-policy update <policy-id> [flags]
```

### Options

```
      --cleanup-drives             Enable cleanup drives for OOB enabled servers
      --disable-embedded-nics      Enable disabling embedded NICs
  -h, --help                       help for update
      --label string               New label for the server cleanup policy
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

