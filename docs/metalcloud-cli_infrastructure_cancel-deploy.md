## metalcloud-cli infrastructure cancel-deploy

Cancel an ongoing infrastructure deployment

### Synopsis

Cancel an ongoing deployment for the specified infrastructure.

This command stops a deployment that is currently in progress, allowing you to
make changes or fix issues before redeploying.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # Cancel deployment by infrastructure ID
  metalcloud-cli infrastructure cancel-deploy 123

  # Cancel deployment by label
  metalcloud-cli infrastructure cancel-deploy "web-cluster"

  # Using the alias
  metalcloud-cli infrastructure cancel my-infrastructure

```
metalcloud-cli infrastructure cancel-deploy infrastructure_id_or_label [flags]
```

### Options

```
  -h, --help   help for cancel-deploy
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

* [metalcloud-cli infrastructure](metalcloud-cli_infrastructure.md)	 - Manage infrastructure resources and configurations

