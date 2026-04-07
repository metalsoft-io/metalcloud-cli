## metalcloud-cli network-device default-secrets batch-create

Create multiple network device default secrets from a CSV file

### Synopsis

Create multiple network device default secrets from a CSV file.

The CSV file must have a header row with the following columns:
  siteId, macAddressOrSerialNumber, secretName, secretValue

Each subsequent row creates one secret.

Required Flags:
  --csv-file    Path to the CSV file

Example CSV file:
  siteId,macAddressOrSerialNumber,secretName,secretValue
  1,AA:BB:CC:DD:EE:FF,admin_password,s3cur3
  1,SN123456,enable_secret,mypass
  2,11:22:33:44:55:66,snmp_community,public

Examples:
  # Batch create secrets from a CSV file
  metalcloud-cli network-device default-secrets batch-create --csv-file secrets.csv

  # Using alias
  metalcloud-cli nd ds batch-create --csv-file /path/to/secrets.csv

```
metalcloud-cli network-device default-secrets batch-create [flags]
```

### Options

```
      --csv-file string   Path to the CSV file
  -h, --help              help for batch-create
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

* [metalcloud-cli network-device default-secrets](metalcloud-cli_network-device_default-secrets.md)	 - Manage network device default secrets

