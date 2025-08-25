## metalcloud-cli infrastructure update

Update infrastructure configuration and metadata

### Synopsis

Update various properties of an infrastructure including its label and custom variables.

This command allows you to modify infrastructure metadata without affecting the deployed
resources. Changes to the infrastructure configuration require a subsequent deploy to
take effect on the actual infrastructure.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure to update
  new_label                   Optional new label for the infrastructure

Flags:
  --custom-variables  JSON string containing custom variables to set on the infrastructure

Examples:
  # Update infrastructure label
  metalcloud-cli infrastructure update 123 "new-cluster-name"

  # Update only custom variables
  metalcloud-cli infrastructure update web-cluster --custom-variables '{"env":"production","version":"1.2.3"}'

  # Update both label and custom variables
  metalcloud-cli infrastructure update 123 "prod-cluster" --custom-variables '{"tier":"production"}'

  # Using the alias
  metalcloud-cli infrastructure edit my-infrastructure new-name

```
metalcloud-cli infrastructure update infrastructure_id_or_label [new_label] [flags]
```

### Options

```
      --custom-variables string   Set of infrastructure custom variables.
  -h, --help                      help for update
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

