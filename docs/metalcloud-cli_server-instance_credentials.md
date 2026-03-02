## metalcloud-cli server-instance credentials

Get login credentials for a server instance

### Synopsis

Get the login credentials for a server instance.

This command retrieves the initial credentials configured for a server instance,
including username, password, and SSH public key if available.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Get credentials for server instance 5678
  metalcloud-cli server-instance credentials 5678

  # Using alias
  metalcloud-cli inst creds 5678

```
metalcloud-cli server-instance credentials <server_instance_id> [flags]
```

### Options

```
  -h, --help   help for credentials
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

