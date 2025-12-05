## metalcloud-cli server-registration-profile create

Create a new server registration profile

### Synopsis

Create a new server registration profile with specified configuration.

The configuration must be provided via a JSON file or pipe and should include
settings for various aspects of server registration.

Required Flags:
  --name            Name for the server registration profile
  --config-source   Source of configuration (path to JSON file or 'pipe' for stdin)

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

Configuration Options:
  - registerCredentials: "user" (keep existing) or "random" (generate new)
  - minimumNumberOfConnectedInterfaces: Minimum required connected interfaces (default: 0)
  - alwaysDiscoverInterfacesWithBDK: Always use BDK for discovery (default: true)
  - enableTpm: Enable TPM during registration (default: true)
  - enableIntelTxt: Enable Intel TXT (default: true)
  - enableSyslogMonitoring: Enable syslog monitoring (default: true)
  - disableTpmAfterRegistration: Disable TPM after registration (default: false)
  - biosProfile: Array of BIOS settings (key-value pairs)
  - defaultVirtualMediaProtocol: HTTPS, NFS, CIFS, etc (default: HTTPS)
  - resetRaidControllers: Reset RAID controllers to factory defaults (default: true)
  - cleanupDrives: Clean up drives during registration (default: true)
  - recreateRaid: Recreate RAID configuration (default: true)
  - disableEmbeddedNics: Disable embedded NICs (default: true)
  - raidOneDrive: RAID level for 1 drive (default: RAID0)
  - raidTwoDrives: RAID level for 2 drives (default: RAID1)
  - raidEvenNumberMoreThanTwoDrives: RAID level for even 4+ drives (default: RAID10)
  - raidOddNumberMoreThanOneDrive: RAID level for odd 3+ drives (default: RAID5)

Examples:
  # Create profile from JSON file
  metalcloud-cli server-registration-profile create \
    --name "production-profile" \
    --config-source /path/to/config.json

  # Create profile from stdin
  cat config.json | metalcloud-cli server-registration-profile create \
    --name "staging-profile" \
    --config-source pipe

  # Create minimal profile
  echo '{"registerCredentials": "random"}' | metalcloud-cli srp create \
    --name "minimal-profile" \
    --config-source pipe

  # Create profile with custom RAID settings
  cat > raid-config.json << EOF
{
  "registerCredentials": "user",
  "resetRaidControllers": true,
  "cleanupDrives": true,
  "recreateRaid": true,
  "raidOneDrive": "RAID0",
  "raidTwoDrives": "RAID1",
  "raidEvenNumberMoreThanTwoDrives": "RAID6",
  "raidOddNumberMoreThanOneDrive": "RAID5"
}
EOF
  metalcloud-cli srp create --name "raid-profile" --config-source raid-config.json

```
metalcloud-cli server-registration-profile create [flags]
```

### Options

```
      --config-source string   Source of the new registration profile configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
      --name string            Name for the server registration profile
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

