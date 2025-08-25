## metalcloud-cli extension-instance list

List all extension instances in an infrastructure

### Synopsis

List all extension instances deployed within a specific infrastructure.

This command displays all extension instances that are currently deployed or configured
within the specified infrastructure. Extension instances represent active deployments
of extensions (workflows, applications, or actions) with their current status,
configuration, and input variables.

The output includes instance details such as:
- Instance ID and label
- Associated extension information
- Current status and state
- Input variables and configuration
- Deployment timestamps

Arguments:
  infrastructure_id_or_label    The unique ID or label of the infrastructure

Examples:
  # List extension instances by infrastructure ID
  metalcloud extension-instance list 12345
  
  # List extension instances by infrastructure label
  metalcloud extension-instance list production-infrastructure
  
  # List instances in staging environment
  metalcloud extension-instance ls staging-env

```
metalcloud-cli extension-instance list infrastructure_id_or_label [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli extension-instance](metalcloud-cli_extension-instance.md)	 - Manage extension instances within infrastructure deployments

