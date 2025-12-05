## metalcloud-cli server-registration-profile get

Get detailed information about a specific server registration profile

### Synopsis

Get detailed information about a specific server registration profile by its ID.

This command retrieves and displays comprehensive information about a server registration
profile including all configuration settings:

Configuration Details:
  - registerCredentials: Password handling mode (user/random)
  - minimumNumberOfConnectedInterfaces: Minimum required connected interfaces
  - alwaysDiscoverInterfacesWithBDK: BDK interface discovery setting
  - enableTpm: TPM enablement
  - enableIntelTxt: Intel TXT enablement
  - enableSyslogMonitoring: Syslog monitoring configuration
  - disableTpmAfterRegistration: Post-registration TPM status
  - biosProfile: Array of BIOS configuration settings
  - defaultVirtualMediaProtocol: Virtual media protocol (HTTPS, etc)
  - resetRaidControllers: RAID controller reset policy
  - cleanupDrives: Drive cleanup policy
  - recreateRaid: RAID recreation policy
  - disableEmbeddedNics: Embedded NIC configuration
  - raidOneDrive: RAID level for single drive (RAID0)
  - raidTwoDrives: RAID level for two drives (RAID1)
  - raidEvenNumberMoreThanTwoDrives: RAID level for even drives 4+ (RAID10)
  - raidOddNumberMoreThanOneDrive: RAID level for odd drives 3+ (RAID5)

Output Format:
  Use global flags to change output format:
  --format=json    JSON output
  --format=csv     CSV output
  --format=yaml    YAML output

Examples:
  # Get profile details by ID
  metalcloud-cli server-registration-profile get 123

  # Get profile in JSON format
  metalcloud-cli server-registration-profile get 123 --format=json

  # Using alias
  metalcloud-cli srp show 123

```
metalcloud-cli server-registration-profile get <profile-id> [flags]
```

### Options

```
  -h, --help   help for get
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

