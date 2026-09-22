## metalcloud-cli server register-production

Register a server that is already running a production workload

### Synopsis

Register a server that is already running a production workload, keeping the
workload in place.

The server is attached to an existing infrastructure instead of being wiped and
re-provisioned. You can provide the configuration either via command-line flags
or by specifying a configuration source using the --config-source flag. The
configuration source can be a path to a JSON file or 'pipe' to read from
standard input.

Required Flags (when not using --config-source):
  --site-id              Site ID where the server is located
  --infrastructure-id    ID of the infrastructure the server belongs to

Optional Flags:
  --config-source        Source of the configuration. Can be 'pipe' or path to a JSON file.
  --os-template-id       ID of the OS template already installed on the server
  --management-address   IPMI/BMC management IP address
  --username             IPMI/BMC username
  --password             IPMI/BMC password
  --serial-number        Server serial number
  --uuid                 Server UUID
  --bmc-mac-address      BMC MAC address
  --vendor               Server vendor
  --model                Server model
  --registration-profile-id  ID of the server registration profile to use

Flag Dependencies:
  --config-source and --site-id are mutually exclusive
  --site-id and --infrastructure-id must be used together

Examples:
  # Register a production server with its interface connections from a JSON file
  metalcloud-cli server register-production --config-source ./production-server.json

  # Register a production server using flags
  metalcloud-cli server register-production --site-id 1 --infrastructure-id 100 \
    --management-address 10.0.0.1 --username admin --password secret

  # Show the expected configuration
  metalcloud-cli server config-example register-production


```
metalcloud-cli server register-production [flags]
```

### Options

```
      --bmc-mac-address string        BMC MAC address.
      --config-source string          Source of the production server configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                          help for register-production
      --infrastructure-id int         ID of the infrastructure the server belongs to.
      --management-address string     IPMI/BMC management IP address.
      --model string                  Server model.
      --os-template-id int            ID of the OS template already installed on the server.
      --password string               IPMI/BMC password.
      --registration-profile-id int   ID of the server registration profile to use.
      --serial-number string          Server serial number.
      --site-id int                   Site ID where the server is located.
      --username string               IPMI/BMC username.
      --uuid string                   Server UUID.
      --vendor string                 Server vendor.
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

