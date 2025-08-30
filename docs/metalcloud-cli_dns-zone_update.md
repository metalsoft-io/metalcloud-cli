## metalcloud-cli dns-zone update

Update DNS zone information

### Synopsis

Update DNS zone information.

This command updates DNS zone configuration using a JSON configuration file or 
piped JSON data. The configuration must be provided via the --config-source flag.

Required Arguments:
  dns_zone_id           The ID of the DNS zone to update

Required Flags:
  --config-source       Source of the DNS zone update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update DNS zone using JSON configuration file
  metalcloud-cli dns-zone update 123 --config-source ./zone-update.json

  # Update DNS zone using piped JSON configuration
  echo '{"description": "Updated description"}' | metalcloud-cli dns-zone update 123 --config-source pipe


```
metalcloud-cli dns-zone update dns_zone_id [flags]
```

### Options

```
      --config-source string   Source of the DNS zone update configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli dns-zone](metalcloud-cli_dns-zone.md)	 - DNS Zone management

