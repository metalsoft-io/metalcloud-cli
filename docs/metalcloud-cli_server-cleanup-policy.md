## metalcloud-cli server-cleanup-policy

Manage server cleanup policies for automated server maintenance

### Synopsis

Manage server cleanup policies that define automated maintenance procedures for servers.

Server cleanup policies control how and when servers are automatically cleaned up,
including configuration of cleanup schedules, retention policies, and maintenance actions.

Available Commands:
  list      List all server cleanup policies
  get       Get details of a specific server cleanup policy
  create    Create a new server cleanup policy
  update    Update an existing server cleanup policy
  delete    Delete a server cleanup policy

Examples:
  # List all server cleanup policies
  metalcloud-cli server-cleanup-policy list

  # Get details of a specific policy
  metalcloud-cli server-cleanup-policy get policy-123

  # Create a new server cleanup policy
  metalcloud-cli server-cleanup-policy create --label "my-policy" --cleanup-drives 1 --recreate-raid 1

  # Update an existing policy
  metalcloud-cli server-cleanup-policy update 123 --label "updated-policy"

  # Delete a policy
  metalcloud-cli server-cleanup-policy delete 123

  # Using short aliases
  metalcloud-cli scp list
  metalcloud-cli srv-cp get policy-123

### Options

```
  -h, --help   help for server-cleanup-policy
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
* [metalcloud-cli server-cleanup-policy create](metalcloud-cli_server-cleanup-policy_create.md)	 - Create a new server cleanup policy
* [metalcloud-cli server-cleanup-policy delete](metalcloud-cli_server-cleanup-policy_delete.md)	 - Delete a server cleanup policy
* [metalcloud-cli server-cleanup-policy get](metalcloud-cli_server-cleanup-policy_get.md)	 - Get detailed information about a specific server cleanup policy
* [metalcloud-cli server-cleanup-policy list](metalcloud-cli_server-cleanup-policy_list.md)	 - List all server cleanup policies
* [metalcloud-cli server-cleanup-policy update](metalcloud-cli_server-cleanup-policy_update.md)	 - Update an existing server cleanup policy

