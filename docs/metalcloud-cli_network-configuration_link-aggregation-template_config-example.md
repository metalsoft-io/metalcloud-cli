## metalcloud-cli network-configuration link-aggregation-template config-example

Generate example configuration template for network device link aggregation configuration template

### Synopsis

Generate an example JSON configuration template that can be used to create
or update network device link aggregation configuration templates.

Preparation and configuration fields need to be base64 encoded when submitted.

Accepted field values:
  action:              create, delete, add-member, remove-member
  aggregationType:     lag, mlag, mlag-peer-link
  networkDeviceDriver: cisco_aci51, nvidia_ufm, nexus9000, cumulus42, arista_eos, dell_s4048, hp5800, hp5900, hp5950, dummy, junos, os_10, sonic_enterprise, vmware_vds, cumulus_linux, brocade, nvidia_dpu, dell_s4000, dell_s6010, junos18

Examples:
  # Display example configuration
  metalcloud-cli network-configuration link-aggregation-template config-example -f json

  # Save example to file
  metalcloud-cli network-configuration link-aggregation-template config-example -f json > template.json

```
metalcloud-cli network-configuration link-aggregation-template config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

* [metalcloud-cli network-configuration link-aggregation-template](metalcloud-cli_network-configuration_link-aggregation-template.md)	 - Manage network devices link aggregation configuration templates

