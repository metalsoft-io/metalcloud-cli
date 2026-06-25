## metalcloud-cli extension-instance credentials

Retrieve the credentials of an extension instance

### Synopsis

Retrieve the credentials (resolved input variables) of an extension instance.

This command returns the credential input variables associated with a deployed
extension instance, such as generated passwords or tokens that the extension
exposes for use.

Arguments:
  extension_instance_id    The unique ID of the extension instance

Examples:
  # Get credentials by instance ID
  metalcloud extension-instance credentials 12345

  # Get credentials using alias
  metalcloud ext-inst get-credentials 12345

```
metalcloud-cli extension-instance credentials extension_instance_id [flags]
```

### Options

```
  -h, --help   help for credentials
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

