## metalcloud-cli fabric link-aggregation

Manage the link aggregations of a fabric

### Synopsis

Manage the link aggregations (LAG / MLAG) of a network fabric.

A link aggregation groups several fabric links into one logical link and
carries the custom variables exposed to the device configuration templates as
link_aggregation_custom_variables.

Commands:
  list, get, create, update, delete, config-example

### Options

```
  -h, --help   help for link-aggregation
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
* [metalcloud-cli fabric link-aggregation config-example](metalcloud-cli_fabric_link-aggregation_config-example.md)	 - Show a link aggregation configuration example
* [metalcloud-cli fabric link-aggregation create](metalcloud-cli_fabric_link-aggregation_create.md)	 - Create a link aggregation on a fabric
* [metalcloud-cli fabric link-aggregation delete](metalcloud-cli_fabric_link-aggregation_delete.md)	 - Delete a link aggregation of a fabric
* [metalcloud-cli fabric link-aggregation get](metalcloud-cli_fabric_link-aggregation_get.md)	 - Get one link aggregation of a fabric
* [metalcloud-cli fabric link-aggregation list](metalcloud-cli_fabric_link-aggregation_list.md)	 - List the link aggregations of a fabric
* [metalcloud-cli fabric link-aggregation update](metalcloud-cli_fabric_link-aggregation_update.md)	 - Update a link aggregation of a fabric

