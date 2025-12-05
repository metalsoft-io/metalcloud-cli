## metalcloud-cli server-registration-profile update

Update an existing server registration profile

### Synopsis

Update an existing server registration profile by its ID.

You can update the profile name, settings, or both. At least one of --name or
--config-source must be provided.

Required Arguments:
  <profile-id>      ID of the server registration profile to update

Optional Flags:
  --name            New name for the server registration profile
  --config-source   Source of configuration updates (path to JSON file or 'pipe')

Note: At least one of --name or --config-source must be specified.

Configuration Structure (JSON):
{
  "registerCredentials": "user|random",
  "minimumNumberOfConnectedInterfaces": 0,
  "alwaysDiscoverInterfacesWithBDK": true,
  "enableTpm": true,
  "enableIntelTxt": true,
  "enableSyslogMonitoring": true,
  "disableTpmAfterRegistration": false,
  "biosProfile": [
    {
      "key": "BootMode",
      "value": "Uefi"
    }
  ],
  "defaultVirtualMediaProtocol": "HTTPS",
  "resetRaidControllers": true,
  "cleanupDrives": true,
  "recreateRaid": true,
  "disableEmbeddedNics": true,
  "raidOneDrive": "RAID0",
  "raidTwoDrives": "RAID1",
  "raidEvenNumberMoreThanTwoDrives": "RAID10",
  "raidOddNumberMoreThanOneDrive": "RAID5"
}

Update Behavior:
  - Only specified fields in the configuration will be updated
  - Unspecified fields will retain their current values
  - Use --name to update only the profile name
  - Use --config-source to update only the settings
  - Use both to update name and settings simultaneously

Examples:
  # Update profile name only
  metalcloud-cli server-registration-profile update 123 \
    --name "updated-profile-name"

  # Update profile settings from file
  metalcloud-cli server-registration-profile update 123 \
    --config-source /path/to/updates.json

  # Update both name and settings
  metalcloud-cli server-registration-profile update 123 \
    --name "new-name" \
    --config-source updates.json

  # Update from stdin
  cat updates.json | metalcloud-cli srp update 123 --config-source pipe

  # Update specific settings only
  echo '{"enableTpm": false, "enableIntelTxt": false}' | \
    metalcloud-cli srp update 123 --config-source pipe

  # Disable drive cleanup and RAID recreation
  cat > updates.json << EOF
{
  "cleanupDrives": false,
  "recreateRaid": false
}
EOF
  metalcloud-cli srp update 123 --config-source updates.json

  # Change RAID configuration strategy
  echo '{"raidTwoDrives": "RAID0", "raidEvenNumberMoreThanTwoDrives": "RAID5"}' | \
    metalcloud-cli server-registration-profile update 123 --config-source pipe

```
metalcloud-cli server-registration-profile update <profile-id> [flags]
```

### Options

```
      --config-source string   Source of the registration profile configuration updates. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
      --name string            New name for the server registration profile
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

* [metalcloud-cli server-registration-profile](metalcloud-cli_server-registration-profile.md)	 - Manage server registration profiles

