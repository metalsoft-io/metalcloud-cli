## metalcloud-cli fabric configure-switches-example

Show an example switch configuration for configure-switches

### Synopsis

Print a commented, ready-to-edit example of the configuration accepted by
'fabric configure-switches'. The output is valid YAML; redirect it to a file or
pipe it straight into the command.

Examples:
  metalcloud-cli fabric configure-switches-example
  metalcloud-cli fabric configure-switches-example > fabric-config.yaml
  metalcloud-cli fabric configure-switches-example | metalcloud-cli fabric configure-switches 5 --config-source pipe --dry-run

```
metalcloud-cli fabric configure-switches-example [flags]
```

### Options

```
  -h, --help   help for configure-switches-example
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

