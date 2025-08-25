## metalcloud-cli vm-pool create

Create a new VM pool from configuration file or pipe

### Synopsis

Create a new virtual machine pool from a JSON configuration file or piped input.

This command creates a VM pool by reading configuration from either a file or standard input.
The configuration must include all required fields and may include optional fields for 
complete setup.

REQUIRED FLAGS:
  --config-source  Source of the VM pool configuration (required)
                   Values: 'pipe' for stdin input, or path to JSON file

CONFIGURATION FIELDS:
  Required:
  - siteId         Site ID where the VM pool will be created
  - managementHost Hostname or IP of the hypervisor management interface
  - managementPort Port for management interface (typically 443 for VMware)
  - name           Name for the VM pool
  - type           VM pool type (e.g., vmware, hyperv, kvm, xen)

  Optional:
  - description    Descriptive text for the VM pool
  - certificate    TLS certificate for secure connections
  - privateKey     Private key corresponding to the certificate
  - username       Username for authentication (alternative to certificates)
  - password       Password for authentication (alternative to certificates)
  - inMaintenance  Set to 1 to create in maintenance mode (default: 0)
  - isExperimental Set to 1 to mark as experimental (default: 0)
  - tags           Array of string tags for categorization
  - options        Additional configuration options specific to the pool type

EXAMPLES:
  # Create from file
  metalcloud-cli vm-pool create --config-source vmpool.json

  # Create from pipe using config example as template
  metalcloud-cli vm-pool config-example | jq '.siteId = 2 | .name = "Production-VMware"' | metalcloud-cli vm-pool create --config-source pipe

  # Create minimal VMware pool from pipe
  echo '{"siteId":1,"managementHost":"vcenter.company.com","managementPort":443,"name":"Test-Pool","type":"vmware"}' | metalcloud-cli vm-pool create --config-source pipe

```
metalcloud-cli vm-pool create [flags]
```

### Options

```
      --config-source string   Source of the new VM pool configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

* [metalcloud-cli vm-pool](metalcloud-cli_vm-pool.md)	 - Manage virtual machine pools and their resources

