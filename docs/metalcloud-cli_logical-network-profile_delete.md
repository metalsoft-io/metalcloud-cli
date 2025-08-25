## metalcloud-cli logical-network-profile delete

Delete a logical network profile

### Synopsis

Delete a logical network profile from the system.

This command permanently removes a logical network profile and all its associated
configuration. The profile must not be in use by any active deployments before
it can be deleted.

Required Arguments:
  logical_network_profile_id    The unique identifier of the profile to delete

Examples:
  # Delete profile by ID
  metalcloud-cli logical-network-profile delete 12345

  # Delete profile using short alias
  metalcloud-cli lnp rm 12345

  # Delete profile using alias
  metalcloud-cli network-profile delete 12345

Warning: This operation is irreversible. Ensure the profile is not in use
by any active infrastructure deployments before deletion.

```
metalcloud-cli logical-network-profile delete logical_network_profile_id [flags]
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

* [metalcloud-cli logical-network-profile](metalcloud-cli_logical-network-profile.md)	 - Manage logical network profiles for network configuration templates

