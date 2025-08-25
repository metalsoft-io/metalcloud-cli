## metalcloud-cli server update-ipmi-credentials

Update server IPMI credentials

### Synopsis

Update server IPMI credentials.

This command updates the IPMI/BMC username and password for the specified server.
The credentials are used for server management operations like power control
and hardware monitoring.

Required Arguments:
  server_id              The ID of the server to update
  username               New IPMI/BMC username
  password               New IPMI/BMC password

Examples:
  # Update IPMI credentials for server with ID 123
  metalcloud-cli server update-ipmi-credentials 123 admin newpassword


```
metalcloud-cli server update-ipmi-credentials server_id username password [flags]
```

### Options

```
  -h, --help   help for update-ipmi-credentials
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

