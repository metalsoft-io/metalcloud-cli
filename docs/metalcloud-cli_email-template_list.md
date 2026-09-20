## metalcloud-cli email-template list

List email templates

### Synopsis

List all email templates defined on the platform.

The table omits the text and html bodies because they are multi-line documents.
Use 'email-template get <name> -f json' to read them.

Optional Flags:
  --filter-name strings   Filter by template name. Repeatable or comma-separated.

Examples:
  metalcloud email-template list
  metalcloud email-templates ls --filter-name infrastructure-deployed

```
metalcloud-cli email-template list [flags]
```

### Options

```
      --filter-name strings   Filter by template name.
  -h, --help                  help for list
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

