## metalcloud-cli server-registration-profile

Manage server registration profiles

### Synopsis

Manage server registration profiles that determine how server registration is performed.

Server registration profiles control how the servers are registered.

Available Commands:
  list      List all server registration profiles
  get       Get details of a specific server registration profile
  create    Create a new server registration profile
  update    Update an existing server registration profile
  delete    Delete a server registration profile

Examples:
  # List all server registration profiles
  metalcloud-cli server-registration-profile list

  # Get details of a specific registration profile
  metalcloud-cli server-registration-profile get srp-123

  # Create a new server registration profile
  metalcloud-cli server-registration-profile create --name "my-profile" --config-source my-profile.json

  # Update an existing registration profile
  metalcloud-cli server-registration-profile update 123 --name "updated-profile"

  # Delete a registration profile
  metalcloud-cli server-registration-profile delete 123

  # Using short aliases
  metalcloud-cli srp list
  metalcloud-cli srv-rp get profile-123

### Options

```
  -h, --help   help for server-registration-profile
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli server-registration-profile create](metalcloud-cli_server-registration-profile_create.md)	 - Create a new server registration profile
* [metalcloud-cli server-registration-profile delete](metalcloud-cli_server-registration-profile_delete.md)	 - Delete a server registration profile
* [metalcloud-cli server-registration-profile get](metalcloud-cli_server-registration-profile_get.md)	 - Get detailed information about a specific server registration profile
* [metalcloud-cli server-registration-profile list](metalcloud-cli_server-registration-profile_list.md)	 - List all server registration profiles
* [metalcloud-cli server-registration-profile update](metalcloud-cli_server-registration-profile_update.md)	 - Update an existing server registration profile

