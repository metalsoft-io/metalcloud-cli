## metalcloud-cli server-default-credentials update

Update existing server default credentials

### Synopsis

Update existing server default credentials from a JSON or YAML configuration.

The configuration may change the default username and password as well as the
rack placement, inventory ID, UUID and registration profile recorded for the
not-yet-registered server. Only the fields present in the configuration are
changed. The site ID, serial number and MAC address that identify the entry
cannot be updated.

Required Arguments:
  credentials_id     The numeric ID of the server default credentials to update

Required Flags:
  --config-source    Source of the update configuration. Can be 'pipe' or path to a JSON file.

Examples:
  # Update credentials from a JSON file
  metalcloud-cli server-default-credentials update 123 --config-source ./credentials.json

  # Update the default password from piped configuration
  echo '{"defaultPassword":"new-secret"}' | metalcloud-cli sdc update 123 --config-source pipe


```
metalcloud-cli server-default-credentials update credentials_id [flags]
```

### Options

```
      --config-source string   Source of the update configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli server-default-credentials](metalcloud-cli_server-default-credentials.md)	 - Manage server default credentials and authentication settings

