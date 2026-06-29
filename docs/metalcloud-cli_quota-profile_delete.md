## metalcloud-cli quota-profile delete

Delete a quota profile

### Synopsis

Delete a quota profile from MetalSoft infrastructure.

This command permanently deletes a quota profile. This action cannot be undone,
so use with caution.

Required Arguments:
  profile_id            The ID of the quota profile to delete

Examples:
  # Delete a quota profile
  metalcloud-cli quota-profile delete example-quota-profile

  # Delete a quota profile using alias
  metalcloud-cli quota-profile rm example-quota-profile


```
metalcloud-cli quota-profile delete profile_id [flags]
```

### Options

```
  -h, --help   help for delete
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

