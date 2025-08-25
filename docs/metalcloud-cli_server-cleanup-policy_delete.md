## metalcloud-cli server-cleanup-policy delete

Delete a server cleanup policy

### Synopsis

Delete a server cleanup policy by its ID.

This command permanently removes a server cleanup policy from the system.
This action cannot be undone, so use with caution.

Arguments:
  policy-id    The unique identifier of the server cleanup policy to delete.
               This must be the numeric ID of the policy.

Required Permissions:
  - server_cleanup_policies:write

Examples:
  # Delete a policy by ID
  metalcloud-cli server-cleanup-policy delete 123

  # Using aliases
  metalcloud-cli scp rm 456
  metalcloud-cli srv-cp remove 789

```
metalcloud-cli server-cleanup-policy delete <policy-id> [flags]
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

* [metalcloud-cli server-cleanup-policy](metalcloud-cli_server-cleanup-policy.md)	 - Manage server cleanup policies for automated server maintenance

