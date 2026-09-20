## metalcloud-cli network-device-controller create

Create a network device controller

### Synopsis

Register a new network device controller.

The controller can be described either with a configuration file (--config-source)
or with individual flags.

Required Flags (one of):
  --config-source string        Source of the controller configuration. Can be 'pipe' or path to a JSON/YAML file.
  --management-address string   Management address of the controller (used together with the flags below)

Required Flags (when not using --config-source):
  --datacenter-name string      Name of the datacenter the controller manages
  --driver string               Controller driver (cisco_ndfc, cisco_aci51, nvidia_ufm, brocade)
  --username string             Management username
  --management-password string  Management password

Optional Flags (when not using --config-source):
  --management-port int         Management port (default 443)
  --site-id int                 ID of the site the controller belongs to
  --identifier-string string    Identifier (hostname) of the controller
  --description string          Free-text description

Examples:
  metalcloud network-device-controller create --config-source controller.json
  metalcloud ndc create --management-address 10.0.0.50 --datacenter-name dc1 \
      --driver cisco_ndfc --username admin --management-password secret

```
metalcloud-cli network-device-controller create [flags]
```

### Options

```
      --config-source string         Source of the new controller configuration. Can be 'pipe' or path to a JSON/YAML file.
      --datacenter-name string       Name of the datacenter the controller manages.
      --description string           Description of the controller.
      --driver string                Controller driver (cisco_ndfc, cisco_aci51, nvidia_ufm, brocade).
  -h, --help                         help for create
      --identifier-string string     Identifier (hostname) of the controller.
      --management-address string    Management address of the controller.
      --management-password string   Management password.
      --management-port int32        Management port. (default 443)
      --site-id int                  ID of the site the controller belongs to.
      --username string              Management username.
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

* [metalcloud-cli network-device-controller](metalcloud-cli_network-device-controller.md)	 - Network device controller management

