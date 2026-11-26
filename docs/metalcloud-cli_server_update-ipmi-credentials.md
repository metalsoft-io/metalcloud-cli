## metalcloud-cli server update-ipmi-credentials

Update server IPMI credentials

### Synopsis

Update server IPMI credentials.

This command updates the IPMI/BMC username and password for the specified server.
The credentials are used for server management operations like power control
and hardware monitoring.

By default the new credentials are only recorded in MetalCloud; the BMC itself
is not changed. Pass --update-on-server to also set them on the BMC (password
rotation). Pass --host to also change the server's management address.

Required Arguments:
  server_id              The ID of the server to update
  username               New IPMI/BMC username
  password               New IPMI/BMC password

Optional Flags:
  --host                 New management (BMC) address of the server
  --update-on-server     Also change the credentials on the BMC itself

Examples:
  # Record new IPMI credentials for server with ID 123
  metalcloud-cli server update-ipmi-credentials 123 admin newpassword

  # Rotate the BMC password on the BMC itself
  metalcloud-cli server update-ipmi-credentials 123 admin newpassword --update-on-server

  # Record new credentials and a new management address
  metalcloud-cli server update-ipmi-credentials 123 admin newpassword --host 10.0.0.15


```
metalcloud-cli server update-ipmi-credentials server_id username password [flags]
```

### Options

```
  -h, --help               help for update-ipmi-credentials
      --host string        New management (BMC) address of the server.
      --update-on-server   Also change the credentials on the BMC itself, not only in MetalCloud.
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API token (msat_...), or a legacy API key (<user id>:<key>), which works only until it expires
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

