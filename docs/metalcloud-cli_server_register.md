## metalcloud-cli server register

Create a server record and optionally register it

### Synopsis

Create a server record in MetalSoft and optionally start its registration.

The --status flag controls what happens after the record is created:
  registering            (default) Start the registration right away. The BMC must be
                         reachable, so --management-address and the BMC credentials
                         are normally needed.
  planned                Only record the server, for example ahead of its delivery.
                         Only --site-id and --serial-number are required; all other
                         details (BMC address, credentials, rack, inventory) are optional
                         and can be filled in later with "server update". Start the
                         registration later with "server register-planned".
  registration_ready     Record the server as ready to be registered, for example
                         when it is to be picked up by ZTP/DHCP discovery.

You can provide the server configuration either via command-line flags or with
--config-source (a JSON/YAML file path or 'pipe'). Run
"metalcloud-cli server config-example register" for an example.

Required Flags (when not using --config-source):
  --site-id                Site ID where the server is located
  --serial-number          Server serial number

Optional Flags:
  --config-source          Source of server configuration (file path or 'pipe')
  --status                 Initial status: registering (default), planned, registration_ready
  --management-address     BMC management IP address
  --username               BMC username
  --password               BMC password
  --model                  Server model
  --vendor                 Server vendor
  --uuid                   Server UUID
  --bmc-mac                BMC MAC address
  --registration-profile-id  Registration profile to use
  --rack-name              Rack name
  --rack-position-lower    Rack position lower unit
  --rack-position-upper    Rack position upper unit
  --inventory-id           Inventory ID

The rack and inventory details are not accepted when the server is created, so they
are applied with an update of the new server right after it is created.

Examples:
  # Create and register a server right away
  metalcloud-cli server register --site-id 1 --serial-number ABC123 --management-address 10.0.0.1 --username admin --password secret

  # Record a planned server with only its serial number
  metalcloud-cli server register --site-id 1 --serial-number ABC123 --status planned

  # Record a planned server with its BMC defaults and rack position
  metalcloud-cli server register --site-id 1 --serial-number ABC123 --status planned \
    --username root --password calvin --rack-name R01 --rack-position-lower 10 --rack-position-upper 11

  # Register using a configuration file
  metalcloud-cli server register --config-source ./server.json

  # Register using piped configuration
  cat server.json | metalcloud-cli server register --config-source pipe


```
metalcloud-cli server register [flags]
```

### Options

```
      --bmc-mac string                BMC MAC address
      --config-source string          Source of the new server configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                          help for register
      --inventory-id string           Inventory ID
      --management-address string     BMC management address
      --model string                  Model
      --password string               BMC password
      --rack-name string              Rack name
      --rack-position-lower string    Rack position lower unit
      --rack-position-upper string    Rack position upper unit
      --registration-profile-id int   Registration profile ID
      --serial-number string          Serial number
      --site-id int                   Site ID
      --status string                 Initial server status: registering (default), planned or registration_ready
      --username string               BMC username
      --uuid string                   Server UUID
      --vendor string                 Vendor
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

