## metalcloud-cli server-registration-profile delete

Delete a server registration profile

### Synopsis

Delete a server registration profile by its ID.

This command permanently removes a server registration profile from the system.
Once deleted, the profile cannot be recovered.

Required Arguments:
  <profile-id>      ID of the server registration profile to delete

Warning:
  - This operation is irreversible
  - Ensure the profile is not currently in use by any servers
  - Servers using this profile may fail to register properly after deletion

Examples:
  # Delete a server registration profile
  metalcloud-cli server-registration-profile delete 123

  # Using aliases
  metalcloud-cli srp rm 123
  metalcloud-cli srv-rp remove 123

  # Delete with confirmation in script
  PROFILE_ID=123
  metalcloud-cli server-registration-profile delete $PROFILE_ID

```
metalcloud-cli server-registration-profile delete <profile-id> [flags]
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

* [metalcloud-cli server-registration-profile](metalcloud-cli_server-registration-profile.md)	 - Manage server registration profiles

