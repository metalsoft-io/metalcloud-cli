## metalcloud-cli quota-profile update

Update quota profile information

### Synopsis

Update quota profile information.

This command updates quota profile configuration using a JSON configuration file
or piped JSON data. The configuration must be provided via the --config-source flag.

Required Arguments:
  profile_id            The ID of the quota profile to update

Required Flags:
  --config-source       Source of the quota profile update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update quota profile using JSON configuration file
  metalcloud-cli quota-profile update example-quota-profile --config-source ./quota-profile-update.json

  # Update quota profile using piped JSON configuration
  echo '{"description": "Updated description"}' | metalcloud-cli quota-profile update example-quota-profile --config-source pipe


```
metalcloud-cli quota-profile update profile_id [flags]
```

### Options

```
      --config-source string   Source of the quota profile update configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

* [metalcloud-cli quota-profile](metalcloud-cli_quota-profile.md)	 - Quota Profile management

