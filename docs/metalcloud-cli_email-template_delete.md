## metalcloud-cli email-template delete

Delete an email template

### Synopsis

Delete an email template. Notifications that rely on it stop being sent until a
template with the same name is created again.

Required Arguments:
  email_template_name   The name of the email template

Examples:
  metalcloud email-template delete infrastructure-deployed
  metalcloud email-template rm infrastructure-deployed

```
metalcloud-cli email-template delete email_template_name [flags]
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

* [metalcloud-cli email-template](metalcloud-cli_email-template.md)	 - Email template management

