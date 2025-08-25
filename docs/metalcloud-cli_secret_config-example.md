## metalcloud-cli secret config-example

Show example secret configuration format

### Synopsis

Display an example JSON configuration structure for creating secrets.

This command outputs a sample configuration that can be used as a template
for creating new secrets. The configuration includes all available fields
with example values.

The configuration format includes:
- name: The secret name (required)
- value: The secret value to encrypt (required)  
- usage: The usage type (optional, defaults to "credential")

Available usage types:
- credential: For storing passwords, API keys, tokens
- configuration: For storing configuration values
- certificate: For storing SSL/TLS certificates
- ssh_key: For storing SSH keys

Examples:
  # Show configuration example
  metalcloud-cli secret config-example

  # Save example to file for editing
  metalcloud-cli secret config-example > my-secret.json

  # Use example as template with custom output format
  metalcloud-cli secret config-example --output json

```
metalcloud-cli secret config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

* [metalcloud-cli secret](metalcloud-cli_secret.md)	 - Manage encrypted secrets for secure credential storage

