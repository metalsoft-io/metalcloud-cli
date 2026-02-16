## metalcloud-cli fabric add-link

Add a network fabric link

### Synopsis

Add a new network fabric link to an existing fabric.

This command creates a new link in the fabric using the configuration provided through
the --config-source flag. The configuration must be a JSON file or piped input containing
the link details such as source and destination network devices and interfaces.

Arguments:
  fabric_id     The ID or label of the fabric to add the link to

Required Flags when using raw configuration:
  --config-source string   Source of the link configuration. Can be 'pipe' for piped input
                          or path to a JSON file containing the link configuration

Required Flags when using individual flags:
  --networkDeviceA string  Identifier string of network device A
  --InterfaceA     string  Name of the interface A
  --networkDeviceB string  Identifier string of network device B
  --InterfaceB     string  Name of the interface B
  --linkType       string  Link type: point-to-point, broadcast

Optional Flags when using individual flags:
  --bgpNumbering         string   inherited, numbered, unnumbered
  --bgpLinkConfiguration string   disabled, active, passive
  --customVariables

Examples:
  # Add link with configuration from file
  metalcloud fabric add-link my-fabric --config-source link-config.json
  
  # Add link with piped configuration
  cat link-config.json | metalcloud fabric add-link 12345 --config-source pipe
  
  # Using alias
  metalcloud fabric create-link my-fabric --config-source link.json

```
metalcloud-cli fabric add-link fabric_id [flags]
```

### Options

```
      --bgp-link-configuration string   BGP configuration (default "disabled")
      --bgp-numbering string            BGP numbering (default "inherited")
      --config-source string            Source of the link configuration. Can be 'pipe' or path to a JSON file.
      --custom-variable stringArray     Custom variable
  -h, --help                            help for add-link
      --interface-a string              Name of the interface A
      --interface-b string              Name of the interface B
      --link-type string                Type of the link
      --mlag-pair                       Set to true if the link on part of MLAG pair
      --network-device-a string         Identifier of the network device A
      --network-device-b string         Identifier of the network device B
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

