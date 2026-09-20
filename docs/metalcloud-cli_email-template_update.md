## metalcloud-cli email-template update

Update an email template

### Synopsis

Update the subject, description, text or html of an email template. Fields absent
from the configuration document are left unchanged. The template name cannot be
changed.

Required Arguments:
  email_template_name      The name of the email template

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud email-template update infrastructure-deployed --config-source update.json
  echo '{"subject":"New subject"}' | metalcloud email-template update infrastructure-deployed --config-source pipe

```
metalcloud-cli email-template update email_template_name [flags]
```

### Options

```
      --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update
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

