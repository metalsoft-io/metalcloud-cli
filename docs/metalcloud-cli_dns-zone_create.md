## metalcloud-cli dns-zone create

Create a new DNS zone

### Synopsis

Create a new DNS zone in MetalSoft.

You can provide the DNS zone configuration either via command-line flags or by 
specifying a configuration source using the --config-source flag. The configuration 
source can be a path to a JSON file or 'pipe' to read from standard input.

If --config-source is not provided, you must specify at least --zone-name, 
--is-default, and --name-servers, along with any other relevant zone details.

Required Flags (when not using --config-source):
  --zone-name           DNS zone name (without terminating dot)
  --is-default          Whether this is the default DNS zone
  --name-servers        List of name servers (comma-separated)

Optional Flags:
  --config-source       Source of DNS zone configuration (JSON file path or 'pipe')
  --description         DNS zone description
  --zone-type           Zone type (master/slave, default: master)
  --soa-email          Email address of DNS zone administrator
  --ttl                TTL (Time to Live) for the DNS zone
  --tags               Tags for the DNS zone (comma-separated)

Examples:
  # Create using command line flags
  metalcloud-cli dns-zone create --zone-name example.com --is-default true --name-servers ns1.example.com,ns2.example.com

  # Create with additional details
  metalcloud-cli dns-zone create --zone-name test.com --is-default false --name-servers ns1.test.com --description "Test zone" --zone-type master --ttl 300

  # Create using JSON configuration file
  metalcloud-cli dns-zone create --config-source ./zone.json

  # Create using piped JSON configuration
  cat zone.json | metalcloud-cli dns-zone create --config-source pipe


```
metalcloud-cli dns-zone create [flags]
```

### Options

```
      --config-source string   Source of the new DNS zone configuration. Can be 'pipe' or path to a JSON file.
      --description string     DNS zone description
  -h, --help                   help for create
      --is-default             Whether this is the default DNS zone
      --name-servers strings   Name servers for the DNS zone
      --soa-email string       Email address of DNS zone administrator
      --tags strings           Tags for the DNS zone
      --ttl int                TTL (Time to Live) for the DNS zone
      --zone-name string       DNS zone name (without terminating dot)
      --zone-type string       Zone type (master/slave) (default "master")
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

