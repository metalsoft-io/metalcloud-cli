## metalcloud-cli vm-pool import-vms

Import VMs from the hypervisor into a VM pool

### Synopsis

Import virtual machines from the hypervisor management system into a VM pool.

This command imports VMs that exist in the hypervisor (e.g., VMware vCenter, Hyper-V)
but are not yet registered in the MetalCloud VM pool. The import configuration can be
specified either via a configuration file/pipe or by using flags.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool

CONFIGURATION VIA FLAGS:
  --vm-names            Comma-separated list of VM names to import
  --infrastructure-id   Infrastructure ID to associate with the imported VMs

CONFIGURATION VIA FILE/PIPE:
  --config-source      Source of the import configuration
                       Values: 'pipe' for stdin input, or path to JSON file

NOTE: Either use --config-source OR both --infrastructure-id and --vm-names flags.

CONFIGURATION FILE FIELDS:
  - vmNames            Array of VM names to import from the hypervisor
  - infrastructureId   Infrastructure ID to associate with the imported VMs

EXAMPLES:
  # Import VMs using flags
  metalcloud-cli vm-pool import-vms 123 --infrastructure-id 456 --vm-names "vm-prod-01,vm-prod-02,vm-test-01"

  # Import VMs using configuration from file
  metalcloud-cli vm-pool import-vms 123 --config-source import-config.json

  # Import VMs using piped configuration
  echo '{"infrastructureId": 456, "vmNames": ["vm-123", "vm-456"]}' | metalcloud-cli vm-pool import-vms 123 --config-source pipe

```
metalcloud-cli vm-pool import-vms vm_pool_id [flags]
```

### Options

```
      --config-source string       Source of the import configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                       help for import-vms
      --infrastructure-id string   Infrastructure ID to associate with the imported VMs
      --vm-names string            Comma-separated list of VM names to import
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

