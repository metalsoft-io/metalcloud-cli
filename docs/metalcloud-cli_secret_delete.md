## metalcloud-cli secret delete

Delete a secret

### Synopsis

Delete a secret by its ID.

This command permanently removes a secret from the system. The secret ID must
be provided as a numeric value. This action cannot be undone.

Arguments:
  secret_id          Numeric ID of the secret to delete (required)

Examples:
  # Delete a secret by ID
  metalcloud-cli secret delete 123

  # Delete a secret with confirmation
  metalcloud-cli secret delete 456 --auto-approve

Note: Be careful when deleting secrets as this action is irreversible.
Make sure the secret is not being used by any infrastructure configurations
before deletion.

```
metalcloud-cli secret delete secret_id [flags]
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

* [metalcloud-cli secret](metalcloud-cli_secret.md)	 - Manage encrypted secrets for secure credential storage

