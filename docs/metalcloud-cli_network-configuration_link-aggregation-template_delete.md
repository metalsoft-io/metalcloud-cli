## metalcloud-cli network-configuration link-aggregation-template delete

Delete a network device link aggregation configuration template from the system

### Synopsis

Delete a network device link aggregation configuration template from the system.

Arguments:
  network_device_link_aggregation_configuration_template_id   The unique identifier of the network device link aggregation configuration template to delete

Examples:
  # Delete network device link aggregation configuration template
  metalcloud-cli network-configuration link-aggregation-template delete 12345

  # Using alias
  metalcloud-cli nc lat rm 12345

```
metalcloud-cli network-configuration link-aggregation-template delete <network_device_link_aggregation_configuration_template_id> [flags]
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

* [metalcloud-cli network-configuration link-aggregation-template](metalcloud-cli_network-configuration_link-aggregation-template.md)	 - Manage network devices link aggregation configuration templates

