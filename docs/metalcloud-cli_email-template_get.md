## metalcloud-cli email-template get

Get email template details

### Synopsis

Get the details of an email template, including its subject and its text and
html bodies.

Required Arguments:
  email_template_name   The name of the email template

Examples:
  metalcloud email-template get infrastructure-deployed
  metalcloud email-template get infrastructure-deployed -f yaml

```
metalcloud-cli email-template get email_template_name [flags]
```

### Options

```
  -h, --help   help for get
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

