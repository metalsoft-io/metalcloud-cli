## metalcloud-cli email-template create

Create an email template

### Synopsis

Create a new email template.

The template is described by a configuration document with the name, subject, text
and html fields (description is optional). Run 'email-template config-example' for
a starting point.

Required Flags:
  --config-source string   Source of the new template configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud email-template create --config-source template.json
  cat template.yaml | metalcloud email-template create --config-source pipe

```
metalcloud-cli email-template create [flags]
```

### Options

```
      --config-source string   Source of the new template configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for create
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

* [metalcloud-cli email-template](metalcloud-cli_email-template.md)	 - Email template management

